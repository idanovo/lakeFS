package staging_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/treeverse/lakefs/pkg/graveler"
	"github.com/treeverse/lakefs/pkg/graveler/staging"
	"github.com/treeverse/lakefs/pkg/kv"
	"github.com/treeverse/lakefs/pkg/kv/kvtest"
	"github.com/treeverse/lakefs/pkg/testutil"
)

func newTestStagingManager(t *testing.T, kvEnabled bool) (context.Context, graveler.StagingManager) {
	t.Helper()
	ctx := context.Background()
	conn, _ := testutil.GetDB(t, databaseURI)
	if kvEnabled {
		store := kvtest.GetStore(ctx, t)
		return ctx, staging.NewManager(kv.StoreMessage{Store: store})
	}
	return ctx, staging.NewDBManager(conn)
}

func TestStagingManager(t *testing.T) {
	kvE := [2]bool{false, true}
	for _, kvEnabled := range kvE {
		dbStr := ""
		if kvEnabled {
			dbStr = "DB"
		}
		t.Run(fmt.Sprintf("Test%sSetGet", dbStr), func(t *testing.T) {
			testSetGet(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sMultiToken", dbStr), func(t *testing.T) {
			testMultiToken(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sDrop", dbStr), func(t *testing.T) {
			testDrop(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sDropByPrefix", dbStr), func(t *testing.T) {
			testDropByPrefix(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sDropPrefixBytes", dbStr), func(t *testing.T) {
			testDropPrefixBytes(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sList", dbStr), func(t *testing.T) {
			testList(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sSeek", dbStr), func(t *testing.T) {
			testSeek(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sNilValue", dbStr), func(t *testing.T) {
			testNilValue(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sNilIdentity", dbStr), func(t *testing.T) {
			testNilIdentity(t, kvEnabled)
		})
		t.Run(fmt.Sprintf("Test%sDeleteAndTombstone", dbStr), func(t *testing.T) {
			testDeleteAndTombstone(t, kvEnabled)
		})
		if kvEnabled {
			t.Run("TestUpdateKV", func(t *testing.T) {
				testUpdate(t)
			})
		}
	}
}

func testUpdate(t *testing.T) {
	ctx, s := newTestStagingManager(t, true)

	key := "a/b/c/my-key-1234"
	testVal := newTestValue("identity1", "value1")
	testVal2 := newTestValue("identity2", "value2")
	// update missing key
	err := s.Update(ctx, "t1", []byte(key), func(value *graveler.Value) (*graveler.Value, error) {
		require.Nil(t, value)
		return testVal, nil
	})
	require.NoError(t, err)

	// update existing key
	err = s.Update(ctx, "t1", []byte(key), func(value *graveler.Value) (*graveler.Value, error) {
		require.Equal(t, testVal, value)
		return testVal2, nil
	})
	require.NoError(t, err)

	// get value
	val, err := s.Get(ctx, "t1", []byte(key))
	require.NoError(t, err)
	require.Equal(t, testVal2, val)

	// update with error
	err = s.Update(ctx, "t1", []byte(key), func(value *graveler.Value) (*graveler.Value, error) {
		require.Equal(t, testVal2, value)
		return testVal, graveler.ErrNotUnique
	})
	require.ErrorIs(t, err, graveler.ErrNotUnique)

	// get value - see it didn't update after an error
	val, err = s.Get(ctx, "t1", []byte(key))
	require.NoError(t, err)
	require.Equal(t, testVal2, val)
}

func testSetGet(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	_, err := s.Get(ctx, "t1", []byte("a/b/c/"))
	if !errors.Is(err, graveler.ErrNotFound) {
		t.Fatalf("error different than expected. expected=%v, got=%v", graveler.ErrNotFound, err)
	}
	value := newTestValue("identity1", "value1")
	err = s.Set(ctx, "t1", []byte("a/b/c/"), value, true)
	testutil.Must(t, err)
	e, err := s.Get(ctx, "t1", []byte("a/b/c/"))
	testutil.Must(t, err)
	if string(e.Identity) != "identity1" {
		t.Errorf("got wrong value. expected=%s, got=%s", "identity1", string(e.Identity))
	}

	if !kvEnabled {
		t.Run("test overwrites", func(t *testing.T) {
			err = s.Set(ctx, "t2", []byte("a/b/c/d"), value, false)
			testutil.Must(t, err)

			err = s.Set(ctx, "t2", []byte("a/b/c/d"), value, false)
			require.ErrorIs(t, graveler.ErrPreconditionFailed, err)
		})
	}
}

func testMultiToken(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	_, err := s.Get(ctx, "t1", []byte("a/b/c/"))
	if !errors.Is(err, graveler.ErrNotFound) {
		t.Fatalf("error different than expected. expected=%v, got=%v", graveler.ErrNotFound, err)
	}
	err = s.Set(ctx, "t1", []byte("a/b/c/"), newTestValue("identity1", "value1"), true)
	testutil.Must(t, err)
	e, err := s.Get(ctx, "t1", []byte("a/b/c/"))
	testutil.Must(t, err)
	if string(e.Identity) != "identity1" {
		t.Errorf("got wrong identity. expected=%s, got=%s", "identity1", string(e.Identity))
	}
	err = s.Set(ctx, "t2", []byte("a/b/c/"), newTestValue("identity2", "value2"), true)
	testutil.Must(t, err)
	e, err = s.Get(ctx, "t1", []byte("a/b/c/"))
	testutil.Must(t, err)
	if string(e.Identity) != "identity1" {
		t.Errorf("got wrong value identity. expected=%s, got=%s", "identity1", string(e.Identity))
	}
	e, err = s.Get(ctx, "t2", []byte("a/b/c/"))
	testutil.Must(t, err)
	if string(e.Identity) != "identity2" {
		t.Errorf("got wrong value identity. expected=%s, got=%s", "identity2", string(e.Identity))
		t.Errorf("got wrong value identity. expected=%s, got=%s", "identity2", string(e.Identity))
	}
}

func testDrop(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	numOfValues := 1400
	for i := 0; i < numOfValues; i++ {
		err := s.Set(ctx, "t1", []byte(fmt.Sprintf("key%04d", i)), newTestValue(fmt.Sprintf("identity%d", i), fmt.Sprintf("value%d", i)), true)
		testutil.Must(t, err)
		err = s.Set(ctx, "t2", []byte(fmt.Sprintf("key%04d", i)), newTestValue(fmt.Sprintf("identity%d", i), fmt.Sprintf("value%d", i)), true)
		testutil.Must(t, err)
	}
	err := s.Drop(ctx, "t1")
	testutil.Must(t, err)
	v, err := s.Get(ctx, "t1", []byte("key0000"))
	if !errors.Is(err, graveler.ErrNotFound) {
		t.Fatalf("after dropping staging area, expected ErrNotFound in Get. got err=%v, got value=%v", err, v)
	}
	it, _ := s.List(ctx, "t1", 0)
	if it.Next() {
		t.Fatal("expected staging area with token t1 to be empty, got non-empty iterator")
	}
	it.Close()
	it, _ = s.List(ctx, "t2", 0)
	count := 0
	for it.Next() {
		if string(it.Value().Data) != fmt.Sprintf("value%d", count) {
			t.Fatalf("unexpected value returned from List at index %d. expected=%s, got=%s", count, fmt.Sprintf("value%d", count), string(it.Value().Data))
		}
		count++
	}
	it.Close()
	if count != numOfValues {
		t.Errorf("got unexpected number of results. expected=%d, got=%d", numOfValues, count)
	}
}

func testDropByPrefix(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	numOfValues := 2400
	for i := 0; i < numOfValues; i++ {
		err := s.Set(ctx, "t1", []byte(fmt.Sprintf("key%04d", i)), newTestValue(fmt.Sprintf("identity%d", i), fmt.Sprintf("value%d", i)), true)
		testutil.Must(t, err)
		err = s.Set(ctx, "t2", []byte(fmt.Sprintf("key%04d", i)), newTestValue(fmt.Sprintf("identity%d", i), fmt.Sprintf("value%d", i)), true)
		testutil.Must(t, err)
	}
	err := s.DropByPrefix(ctx, "t1", []byte("key1"))
	testutil.Must(t, err)
	v, err := s.Get(ctx, "t1", []byte("key1000"))
	if !errors.Is(err, graveler.ErrNotFound) {
		// key1000 starts with the deleted prefix - should have been deleted
		t.Fatalf("after dropping staging area, expected ErrNotFound in Get. got err=%v, got value=%s", err, v)
	}
	_, err = s.Get(ctx, "t1", []byte("key0000"))
	// key0000 does not start with the deleted prefix - should be returned
	testutil.Must(t, err)
	it, _ := s.List(ctx, "t1", 0)
	count := 0
	for it.Next() {
		count++
	}
	it.Close()
	if count != numOfValues-1000 {
		t.Errorf("got unexpected number of results after drop. expected=%d, got=%d", numOfValues-1000, count)
	}
	it, _ = s.List(ctx, "t2", 0)
	count = 0
	for it.Next() {
		count++
	}
	it.Close()
	if count != numOfValues {
		t.Errorf("got unexpected number of results. expected=%d, got=%d", numOfValues, count)
	}
}

func testDropPrefixBytes(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	tests := map[string]struct {
		keys                    []graveler.Key
		prefix                  graveler.Key
		expectedLengthAfterDrop int
	}{
		"prefix with all bytes=MaxUint8": {
			keys:                    []graveler.Key{{255, 255, 254, 254}, {255, 255, 254, 255}, {255, 255, 255, 253}, {255, 255, 255, 254}, {255, 255, 255, 255}},
			prefix:                  graveler.Key{255, 255, 255},
			expectedLengthAfterDrop: 2,
		},
		"all zero prefix": {
			keys:                    []graveler.Key{{0, 0, 0, 0}, {0, 0, 0, 255}, {0, 0, 1, 0}, {0, 0, 1, 1}},
			prefix:                  graveler.Key{0, 0, 0},
			expectedLengthAfterDrop: 2,
		},
		"prefix common to all keys": {
			keys:                    []graveler.Key{{0, 0, 0, 0}, {0, 0, 0, 255}, {0, 0, 1, 0}, {0, 0, 1, 1}},
			prefix:                  graveler.Key{0, 0},
			expectedLengthAfterDrop: 0,
		},
		"axUint8 in keys - prefix length 1": {
			keys:                    []graveler.Key{{1, 0, 0, 0}, {1, 0, 0, 255}, {1, 0, 255, 255}, {1, 255, 255, 255}},
			prefix:                  graveler.Key{1},
			expectedLengthAfterDrop: 0,
		},
		"MaxUint8 in keys - prefix length 2": {
			keys:                    []graveler.Key{{1, 0, 0, 0}, {1, 0, 0, 255}, {1, 0, 255, 255}, {1, 255, 255, 255}},
			prefix:                  graveler.Key{1, 0},
			expectedLengthAfterDrop: 1,
		},
		"MaxUint8 in keys - prefix length 3": {
			keys:                    []graveler.Key{{1, 0, 0, 0}, {1, 0, 0, 255}, {1, 0, 255, 255}, {1, 255, 255, 255}},
			prefix:                  graveler.Key{1, 0, 0},
			expectedLengthAfterDrop: 2,
		},
		"MaxUint8 in keys - prefix length 4": {
			keys:                    []graveler.Key{{1, 0, 0, 0}, {1, 0, 0, 255}, {1, 0, 255, 255}, {1, 255, 255, 255}},
			prefix:                  graveler.Key{1, 0, 0, 0},
			expectedLengthAfterDrop: 3,
		},
		"multi-length keys - prefix length 3": {
			keys:                    []graveler.Key{{1, 0}, {1, 1}, {1, 0, 1}, {1, 1, 1}, {1, 1, 1, 255}, {1, 1, 255, 1}, {1, 1, 1, 1, 1}, {1, 1, 1, 255, 1, 1, 1, 1, 1, 1, 1, 1}},
			prefix:                  graveler.Key{1, 1, 1},
			expectedLengthAfterDrop: 4,
		},
		"multi-length keys - prefix length 4": {
			keys:                    []graveler.Key{{1, 0}, {1, 1}, {1, 0, 1}, {1, 1, 1}, {1, 1, 1, 255}, {1, 1, 255, 1}, {1, 1, 1, 1, 1}, {1, 1, 1, 255, 1, 1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 255, 255, 255, 255}},
			prefix:                  graveler.Key{1, 1, 1, 1},
			expectedLengthAfterDrop: 7,
		},
		"empty prefix": {
			keys:                    []graveler.Key{{1, 0}, {1, 1}, {1, 0, 1}, {1, 1, 1}, {1, 1, 1, 255}, {1, 1, 255, 1}, {1, 1, 1, 1, 1}, {1, 1, 1, 255, 1, 1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 255, 255, 255, 255}, {2, 0}},
			prefix:                  graveler.Key{},
			expectedLengthAfterDrop: 0,
		},
		"multi-length keys - prefix with MaxUint 8": {
			keys:                    []graveler.Key{{1, 0}, {1, 1}, {1, 0, 1}, {1, 1, 1}, {1, 1, 1, 255}, {1, 1, 255, 1}, {1, 1, 1, 1, 1}, {1, 1, 1, 255, 1, 1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 255, 255, 255, 255}, {2, 0}},
			prefix:                  graveler.Key{0, 255, 255, 255},
			expectedLengthAfterDrop: 10,
		},
		"multi-length keys - prefix with MaxUint 8 - prefix length 2": {
			keys:                    []graveler.Key{{1, 0}, {1, 1}, {1, 0, 1}, {1, 1, 1}, {1, 1, 1, 255}, {1, 1, 255, 1}, {1, 1, 1, 1, 1}, {1, 1, 1, 255, 1, 1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 255, 255, 255, 255}, {2, 0}},
			prefix:                  graveler.Key{1, 255},
			expectedLengthAfterDrop: 10,
		},
		"multi-length keys - prefix with MaxUint 8 - prefix length 3": {
			keys:                    []graveler.Key{{1, 254, 255, 255}, {1, 255}, {1, 255, 255}, {1, 255, 255, 255}, {2, 255}},
			prefix:                  graveler.Key{1, 255, 255},
			expectedLengthAfterDrop: 3,
		},
	}
	for name, tst := range tests {
		st := graveler.StagingToken(fmt.Sprintf("t_%s", name))
		t.Run(name, func(t *testing.T) {
			for _, k := range tst.keys {
				err := s.Set(ctx, st, k, &graveler.Value{
					Identity: []byte{0, 0, 0, 0, 0, 0},
					Data:     []byte{0, 0, 0, 0, 0, 0},
				}, true)
				testutil.Must(t, err)
			}
			err := s.DropByPrefix(ctx, st, tst.prefix)
			testutil.Must(t, err)
			it, err := s.List(ctx, st, 0)
			testutil.Must(t, err)
			count := 0
			for it.Next() {
				count++
			}
			if count != tst.expectedLengthAfterDrop {
				t.Fatalf("unexpected number of values after drop. expected=%d, got=%d", tst.expectedLengthAfterDrop, count)
			}
			if it.Err() != nil {
				t.Fatalf("got unexpected error: %v", it.Err())
			}
			it.Close()
		})
	}
}

func testList(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	for _, numOfValues := range []int{1, 100, 1000, 1500, 2500} {
		token := graveler.StagingToken(fmt.Sprintf("t_%d", numOfValues))
		for i := 0; i < numOfValues; i++ {
			err := s.Set(ctx, token, []byte(fmt.Sprintf("key%04d", i)), newTestValue(fmt.Sprintf("identity%d", i), fmt.Sprintf("value%d", i)), true)
			testutil.Must(t, err)
		}
		res := make([]*graveler.ValueRecord, 0, numOfValues)
		it, _ := s.List(ctx, token, 0)
		for it.Next() {
			res = append(res, it.Value())
		}
		if it.Err() != nil {
			t.Fatalf("got unexpected error from list: %v", it.Err())
		}
		it.Close()
		if len(res) != numOfValues {
			t.Errorf("got unexpected number of results. expected=%d, got=%d", numOfValues, len(res))
		}
		for i, e := range res {
			if !bytes.Equal(e.Key, []byte(fmt.Sprintf("key%04d", i))) {
				t.Fatalf("got unexpected key from List at index %d: expected: key%04d, got: %s", i, i, string(e.Key))
			}
			if string(e.Data) != fmt.Sprintf("value%d", i) {
				t.Fatalf("unexpected value returned from List at index %d. expected=%s, got=%s", i, fmt.Sprintf("value%d", i), string(e.Data))
			}
		}
	}
}

func testSeek(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	numOfValues := 100
	for i := 0; i < numOfValues; i++ {
		err := s.Set(ctx, "t1", []byte(fmt.Sprintf("key%04d", i)), newTestValue("identity1", "value1"), true)
		testutil.Must(t, err)
	}
	it, _ := s.List(ctx, "t1", 0)
	defer it.Close()
	if it.SeekGE([]byte("key0050")); !it.Next() {
		t.Fatal("iterator seek expected to return true, got false")
	}
	expected := "key0050"
	if !bytes.Equal(it.Value().Key, []byte("key0050")) {
		t.Fatalf("got unexpected key after iterator seek. expected=%s, got=%s", expected, string(it.Value().Key))
	}
	if !it.Next() {
		t.Fatal("iterator next expected to return true, got false")
	}
	expected = "key0051"
	if !bytes.Equal(it.Value().Key, []byte(expected)) {
		t.Fatalf("got unexpected key after iterator seek. expected=%s, got=%s", expected, string(it.Value().Key))
	}
	if it.SeekGE([]byte("key1000")); it.Next() {
		t.Fatal("iterator seek expected to return false, got true")
	}
	if it.SeekGE([]byte("key0060a")); !it.Next() {
		t.Fatal("iterator seek expected to return true, got false")
	}
	expected = "key0061"
	if !bytes.Equal(it.Value().Key, []byte(expected)) {
		t.Fatalf("got unexpected key after iterator seek. expected=%s, got=%s", expected, string(it.Value().Key))
	}
	if !it.Next() {
		t.Fatal("iterator next expected to return true, got false")
	}
}

func testNilValue(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	err := s.Set(ctx, "t1", []byte("key1"), nil, true)
	testutil.Must(t, err)
	err = s.Set(ctx, "t1", []byte("key2"), newTestValue("identity2", "value2"), true)
	testutil.Must(t, err)
	e, err := s.Get(ctx, "t1", []byte("key1"))
	testutil.Must(t, err)
	if e != nil {
		t.Errorf("got unexpected value. expected=nil, got=%s", e)
	}
	it, err := s.List(ctx, "t1", 0)
	testutil.Must(t, err)
	defer it.Close()
	if !it.Next() {
		t.Fatalf("expected to get key from list")
	}

	expected := "key1"
	if !bytes.Equal(it.Value().Key, []byte(expected)) {
		t.Errorf("got unexpected key. expected=key1, got=%s", it.Value().Key)
	}
	if it.Value().Value != nil {
		t.Errorf("got unexpected value. expected=nil, got=%s", it.Value().Value)
	}

	if !it.Next() {
		t.Fatalf("expected to get key from list")
	}
	e = it.Value().Value
	if string(e.Identity) != "identity2" {
		t.Errorf("got wrong identity. expected=%s, got=%s", "identity2", string(e.Identity))
	}
}

func testNilIdentity(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	err := s.Set(ctx, "t1", []byte("key1"), newTestValue("identity1", "value1"), true)
	testutil.Must(t, err)
	err = s.Set(ctx, "t1", []byte("key1"), &graveler.Value{
		Identity: nil,
		Data:     []byte("value1"),
	}, true)
	if !errors.Is(err, graveler.ErrInvalidValue) {
		t.Fatalf("got unexpected error. expected=%v, got=%v", graveler.ErrInvalidValue, err)
	}
	e, err := s.Get(ctx, "t1", []byte("key1"))
	testutil.Must(t, err)
	if string(e.Identity) != "identity1" {
		t.Errorf("got wrong identity. expected=%s, got=%s", "identity1", string(e.Identity))
	}
}

func testDeleteAndTombstone(t *testing.T, kvEnabled bool) {
	ctx, s := newTestStagingManager(t, kvEnabled)
	_, err := s.Get(ctx, "t1", []byte("key1"))
	if !errors.Is(err, graveler.ErrNotFound) {
		t.Fatalf("error different than expected. expected=%v, got=%v", graveler.ErrNotFound, err)
	}
	tombstoneValues := []*graveler.Value{
		{
			Identity: []byte("identity1"),
			Data:     make([]byte, 0),
		},
		{
			Identity: []byte("identity1"),
			Data:     nil,
		},
	}
	for _, val := range tombstoneValues {
		err = s.Set(ctx, "t1", []byte("key1"), val, true)
		testutil.Must(t, err)
		e, err := s.Get(ctx, "t1", []byte("key1"))
		testutil.Must(t, err)
		if len(e.Data) != 0 {
			t.Fatalf("expected empty data, got: %v", e.Data)
		}
		if string(e.Identity) != "identity1" {
			t.Fatalf("got unexpected value identity. expected=%s, got=%s", "identity1", string(e.Identity))
		}
		it, err := s.List(ctx, "t1", 0)
		testutil.Must(t, err)
		if !it.Next() {
			t.Fatalf("expected to get key from list")
		}
		if it.Err() != nil {
			t.Fatalf("unexpected error from iterator: %v", it.Err())
		}
		if len(it.Value().Value.Data) != 0 {
			t.Fatalf("expected empty value data from iterator, got: %v", it.Value().Value.Data)
		}
		it.Close()
	}
	err = s.Set(ctx, "t1", []byte("key1"), newTestValue("identity3", "value3"), true)
	testutil.Must(t, err)
	e, err := s.Get(ctx, "t1", []byte("key1"))
	testutil.Must(t, err)
	if string(e.Identity) != "identity3" {
		t.Fatalf("got unexpected value identity. expected=%s, got=%s", "identity3", string(e.Identity))
	}
	err = s.DropKey(ctx, "t1", []byte("key1"))
	testutil.Must(t, err)
	_, err = s.Get(ctx, "t1", []byte("key1"))
	if !errors.Is(err, graveler.ErrNotFound) {
		t.Fatalf("error different than expected. expected=%v, got=%v", graveler.ErrNotFound, err)
	}
}

func newTestValue(identity, data string) *graveler.Value {
	return &graveler.Value{
		Identity: []byte(identity),
		Data:     []byte(data),
	}
}
