package testutil

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/treeverse/lakefs/pkg/auth"
	"github.com/treeverse/lakefs/pkg/auth/model"
	"github.com/treeverse/lakefs/pkg/block"
	"github.com/treeverse/lakefs/pkg/catalog"
	"github.com/treeverse/lakefs/pkg/config"
	"github.com/treeverse/lakefs/pkg/gateway"
	"github.com/treeverse/lakefs/pkg/gateway/multiparts"
	"github.com/treeverse/lakefs/pkg/kv"
	"github.com/treeverse/lakefs/pkg/kv/kvtest"
	_ "github.com/treeverse/lakefs/pkg/kv/mem"
	kvparams "github.com/treeverse/lakefs/pkg/kv/params"
	"github.com/treeverse/lakefs/pkg/logging"
	"github.com/treeverse/lakefs/pkg/stats"
	"github.com/treeverse/lakefs/pkg/testutil"
)

type Dependencies struct {
	blocks  block.Adapter
	auth    *FakeAuthService
	catalog *catalog.Catalog
}

func GetBasicHandler(t *testing.T, authService *FakeAuthService, databaseURI string, repoName string, kvEnabled bool) (http.Handler, *Dependencies) {
	ctx := context.Background()
	conn, _ := testutil.GetDB(t, databaseURI)
	idTranslator := &testutil.UploadIDTranslator{
		TransMap:   make(map[string]string),
		ExpectedID: "",
		T:          t,
	}
	viper.Set(config.BlockstoreTypeKey, block.BlockstoreTypeMem)
	// Disable KV by default (used for determining KV state by certain packages such as catalog)
	viper.Set("database.kv_enabled", false)

	var (
		multipartsTracker multiparts.Tracker
		storeMessage      *kv.StoreMessage
	)
	if kvEnabled {
		store := kvtest.MakeStoreByName("mem", kvparams.KV{})(t, context.Background())
		defer store.Close()
		storeMessage = &kv.StoreMessage{Store: store}
		multipartsTracker = multiparts.NewTracker(*storeMessage)
		viper.Set("database.kv_enabled", true)
	} else {
		multipartsTracker = multiparts.NewDBTracker(conn)
	}

	blockstoreType, _ := os.LookupEnv(testutil.EnvKeyUseBlockAdapter)
	blockAdapter := testutil.NewBlockAdapterByType(t, idTranslator, blockstoreType)

	conf, err := config.NewConfig()
	testutil.MustDo(t, "config", err)

	c, err := catalog.New(ctx, catalog.Config{
		Config:  conf,
		DB:      conn,
		KVStore: storeMessage,
	})
	testutil.MustDo(t, "build catalog", err)
	t.Cleanup(func() {
		_ = c.Close()
	})

	storageNamespace := os.Getenv("USE_STORAGE_NAMESPACE")
	if storageNamespace == "" {
		storageNamespace = "replay"
	}

	_, err = c.CreateRepository(ctx, repoName, storageNamespace, "main")
	testutil.Must(t, err)

	handler := gateway.NewHandler(authService.Region, c, multipartsTracker, blockAdapter, authService, []string{authService.BareDomain}, &mockCollector{}, nil, config.DefaultAuditLogLevel, true)

	return handler, &Dependencies{
		blocks:  blockAdapter,
		auth:    authService,
		catalog: c,
	}
}

type FakeAuthService struct {
	BareDomain      string `json:"bare_domain"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"access_secret_key"`
	UserID          string `json:"user_id"`
	Region          string `json:"region"`
}

func (m *FakeAuthService) GetCredentials(_ context.Context, accessKey string) (*model.Credential, error) {
	if accessKey != m.AccessKeyID {
		logging.Default().Fatal("access key in recording different than configuration")
	}
	aCred := new(model.Credential)
	aCred.AccessKeyID = accessKey
	aCred.SecretAccessKey = m.SecretAccessKey
	aCred.Username = m.UserID
	return aCred, nil
}

func (m *FakeAuthService) GetUser(_ context.Context, _ string) (*model.User, error) {
	return &model.User{
		CreatedAt: time.Now(),
		Username:  "user"}, nil
}

func (m *FakeAuthService) Authorize(_ context.Context, _ *auth.AuthorizationRequest) (*auth.AuthorizationResponse, error) {
	return &auth.AuthorizationResponse{Allowed: true}, nil
}

type mockCollector struct{}

func (m *mockCollector) CollectMetadata(*stats.Metadata) {}

func (m *mockCollector) CollectEvent(string, string) {}

func (m *mockCollector) SetInstallationID(string) {}

func (m *mockCollector) Close() {}
