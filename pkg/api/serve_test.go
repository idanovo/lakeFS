package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/treeverse/lakefs/pkg/auth/email"
	"github.com/treeverse/lakefs/pkg/kv/kvtest"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/spf13/viper"
	"github.com/treeverse/lakefs/pkg/actions"
	"github.com/treeverse/lakefs/pkg/api"
	"github.com/treeverse/lakefs/pkg/auth"
	"github.com/treeverse/lakefs/pkg/auth/crypt"
	authmodel "github.com/treeverse/lakefs/pkg/auth/model"
	authparams "github.com/treeverse/lakefs/pkg/auth/params"
	"github.com/treeverse/lakefs/pkg/block"
	"github.com/treeverse/lakefs/pkg/catalog"
	"github.com/treeverse/lakefs/pkg/config"
	"github.com/treeverse/lakefs/pkg/db"
	dbparams "github.com/treeverse/lakefs/pkg/db/params"
	"github.com/treeverse/lakefs/pkg/ingest/store"
	"github.com/treeverse/lakefs/pkg/kv"
	"github.com/treeverse/lakefs/pkg/logging"
	"github.com/treeverse/lakefs/pkg/stats"
	"github.com/treeverse/lakefs/pkg/templater"
	"github.com/treeverse/lakefs/pkg/testutil"
	"github.com/treeverse/lakefs/pkg/version"
	"github.com/treeverse/lakefs/templates"
)

const (
	ServerTimeout = 30 * time.Second
)

type dependencies struct {
	blocks      block.Adapter
	catalog     catalog.Interface
	authService *auth.Service
	collector   *nullCollector
}

type nullCollector struct {
	metadata []*stats.Metadata
}

func (m *nullCollector) CollectMetadata(metadata *stats.Metadata) {
	m.metadata = append(m.metadata, metadata)
}

func (m *nullCollector) CollectEvent(_, _ string) {}

func (m *nullCollector) SetInstallationID(_ string) {}

func (m *nullCollector) Close() {}

func createDefaultAdminUser(t testing.TB, clt api.ClientWithResponsesInterface) *authmodel.BaseCredential {
	t.Helper()
	res, err := clt.SetupWithResponse(context.Background(), api.SetupJSONRequestBody{
		Username: "admin",
	})
	testutil.Must(t, err)
	if res.JSON200 == nil {
		t.Fatal("Failed run setup env", res.HTTPResponse.StatusCode, res.HTTPResponse.Status)
	}
	return &authmodel.BaseCredential{
		IssuedDate:      time.Unix(res.JSON200.CreationDate, 0),
		AccessKeyID:     res.JSON200.AccessKeyId,
		SecretAccessKey: res.JSON200.SecretAccessKey,
	}
}

func setupHandlerWithWalkerFactory(t testing.TB, factory catalog.WalkerFactory, kvEnabled bool, opts ...testutil.GetDBOption) (http.Handler, *dependencies) {
	t.Helper()
	ctx := context.Background()
	conn, handlerDatabaseURI := testutil.GetDB(t, databaseURI, opts...)
	viper.Set(config.BlockstoreTypeKey, block.BlockstoreTypeMem)
	viper.Set("database.kv_enabled", kvEnabled)

	collector := &nullCollector{}

	// wire actions
	var (
		actionsStore   actions.Store
		idGen          actions.IDGenerator
		authService    auth.Service
		meta           auth.MetadataManager
		kvStoreMessage *kv.StoreMessage
	)

	cfg, err := config.NewConfig()
	testutil.MustDo(t, "config", err)
	if kvEnabled {
		kvStore := kvtest.GetStore(ctx, t)
		kvStoreMessage = &kv.StoreMessage{Store: kvStore}
		actionsStore = actions.NewActionsKVStore(*kvStoreMessage)
		idGen = &actions.DecreasingIDGenerator{}
		authService = auth.NewKVAuthService(kvStoreMessage, crypt.NewSecretStore([]byte("some secret")), nil, authparams.ServiceCache{
			Enabled: false,
		}, logging.Default())
		meta = auth.NewKVMetadataManager("serve_test", cfg.GetFixedInstallationID(), cfg.GetDatabaseParams().Type, kvStore)
		viper.Set("database.kv_enabled", true)
	} else {
		actionsStore = actions.NewActionsDBStore(conn)
		idGen = &actions.IncreasingIDGenerator{}
		authService = auth.NewDBAuthService(conn, crypt.NewSecretStore([]byte("some secret")), nil, authparams.ServiceCache{
			Enabled: false,
		}, logging.Default())
		meta = auth.NewDBMetadataManager("serve_test", cfg.GetFixedInstallationID(), conn)
	}

	// Do not validate invalid config (missing required fields).
	c, err := catalog.New(ctx, catalog.Config{
		Config:        cfg,
		DB:            conn,
		KVStore:       kvStoreMessage,
		WalkerFactory: factory,
	})
	testutil.MustDo(t, "build catalog", err)

	actionsService := actions.NewService(
		ctx,
		actionsStore,
		catalog.NewActionsSource(c),
		catalog.NewActionsOutputWriter(c.BlockAdapter),
		idGen,
		collector,
		true,
	)

	c.SetHooksHandler(actionsService)

	authenticator := auth.NewBuiltinAuthenticator(authService)
	migrator := db.NewDatabaseMigrator(dbparams.Database{ConnectionString: handlerDatabaseURI})

	t.Cleanup(func() {
		actionsService.Stop()
		_ = c.Close()
	})

	auditChecker := version.NewDefaultAuditChecker(cfg.GetSecurityAuditCheckURL())
	emailParams, _ := cfg.GetEmailParams()
	emailer, err := email.NewEmailer(emailParams)
	templater := templater.NewService(templates.Content, cfg, authService)

	testutil.Must(t, err)
	handler := api.Serve(cfg, c, authenticator, authenticator, authService, c.BlockAdapter, meta, migrator, collector, nil, actionsService, auditChecker, logging.Default(), emailer, templater, nil, nil, nil, nil)

	return handler, &dependencies{
		blocks:      c.BlockAdapter,
		authService: &authService,
		catalog:     c,
		collector:   collector,
	}
}

func setupHandler(t testing.TB, kvEnabled bool, opts ...testutil.GetDBOption) (http.Handler, *dependencies) {
	return setupHandlerWithWalkerFactory(t, store.NewFactory(nil), kvEnabled, opts...)
}

func setupClientByEndpoint(t testing.TB, endpointURL string, accessKeyID, secretAccessKey string) api.ClientWithResponsesInterface {
	t.Helper()

	var opts []api.ClientOption
	if accessKeyID != "" {
		basicAuthProvider, err := securityprovider.NewSecurityProviderBasicAuth(accessKeyID, secretAccessKey)
		if err != nil {
			t.Fatal("basic auth security provider", err)
		}
		opts = append(opts, api.WithRequestEditorFn(basicAuthProvider.Intercept))
	}
	clt, err := api.NewClientWithResponses(endpointURL+api.BaseURL, opts...)
	if err != nil {
		t.Fatal("failed to create lakefs api client:", err)
	}
	return clt
}

func setupServer(t testing.TB, handler http.Handler) *httptest.Server {
	t.Helper()
	if shouldUseServerTimeout() {
		handler = http.TimeoutHandler(handler, ServerTimeout, `{"error": "timeout"}`)
	}
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func shouldUseServerTimeout() bool {
	withServerTimeoutVal := os.Getenv("TEST_WITH_SERVER_TIMEOUT")
	if withServerTimeoutVal == "" {
		return true // default
	}
	withServerTimeout, err := strconv.ParseBool(withServerTimeoutVal)
	if err != nil {
		panic(fmt.Errorf("invalid TEST_WITH_SERVER_TIMEOUT value: %w", err))
	}
	return withServerTimeout
}

func setupClientWithAdmin(t testing.TB, kvEnabled bool, opts ...testutil.GetDBOption) (api.ClientWithResponsesInterface, *dependencies) {
	t.Helper()
	return setupClientWithAdminAndWalkerFactory(t, store.NewFactory(nil), kvEnabled, opts...)
}

func setupClientWithAdminAndWalkerFactory(t testing.TB, factory catalog.WalkerFactory, kvEnabled bool, opts ...testutil.GetDBOption) (api.ClientWithResponsesInterface, *dependencies) {
	t.Helper()
	handler, deps := setupHandlerWithWalkerFactory(t, factory, kvEnabled, opts...)
	server := setupServer(t, handler)
	clt := setupClientByEndpoint(t, server.URL, "", "")
	cred := createDefaultAdminUser(t, clt)
	clt = setupClientByEndpoint(t, server.URL, cred.AccessKeyID, cred.SecretAccessKey)
	return clt, deps
}

func TestInvalidRoute(t *testing.T) {
	handler, _ := setupHandler(t, false)
	server := setupServer(t, handler)
	clt := setupClientByEndpoint(t, server.URL, "", "")
	cred := createDefaultAdminUser(t, clt)

	// setup client with invalid endpoint base url
	basicAuthProvider, err := securityprovider.NewSecurityProviderBasicAuth(cred.AccessKeyID, cred.SecretAccessKey)
	if err != nil {
		t.Fatal("basic auth security provider", err)
	}
	clt, err = api.NewClientWithResponses(server.URL+api.BaseURL+"//", api.WithRequestEditorFn(basicAuthProvider.Intercept))
	if err != nil {
		t.Fatal("failed to create api client:", err)
	}

	ctx := context.Background()
	resp, err := clt.ListRepositoriesWithResponse(ctx, &api.ListRepositoriesParams{})
	if err != nil {
		t.Fatalf("failed to get lakefs server version")
	}
	if resp.JSONDefault == nil {
		t.Fatalf("client api call expected default error, got nil")
	}
	expectedErrMsg := api.ErrInvalidAPIEndpoint.Error()
	errMsg := resp.JSONDefault.Message
	if errMsg != expectedErrMsg {
		t.Fatalf("client response error message: %s, expected: %s", errMsg, expectedErrMsg)
	}
}
