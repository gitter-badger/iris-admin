package cache

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/server/zap_server"
)

func TestSetCacheString(t *testing.T) {
	defer Remove()
	defer zap_server.Remove()
	redisPwd := os.Getenv("redisPwd")
	CONFIG.Password = redisPwd
	t.Run("test set cache string", func(t *testing.T) {
		key := "test_set_cache"
		want := "test_set_cache_value"
		err := SetCache(key, want, time.Duration(time.Second*3))
		if err != nil {
			t.Errorf("set cache get error %v\n", err)
		}
		get, err := GetCacheString(key)
		if err != nil {
			t.Errorf("set cache get error %v\n", err)
		}
		if get != want {
			t.Errorf("set cache want [%s] but get [%s]\n", want, get)
		}
		time.Sleep(time.Second * 5)
		_, err = GetCacheString(key)
		if err == nil {
			t.Error("set cache want error but get nil\n")
		}
	})
}

func TestSetCacheUint(t *testing.T) {
	defer Remove()
	defer zap_server.Remove()
	redisPwd := os.Getenv("redisPwd")
	CONFIG.Password = redisPwd
	t.Run("test set cache bytes", func(t *testing.T) {
		key := "test_set_cache"
		var want uint64 = 123
		err := SetCache(key, want, time.Duration(time.Second*3))
		if err != nil {
			t.Errorf("set cache get error %v\n", err)
		}
		get, err := GetCacheUint(key)
		if err != nil {
			t.Errorf("set cache get error %v\n", err)
		}
		if get != want {
			t.Errorf("set cache want [%d] but get [%d]\n", want, get)
		}
		time.Sleep(time.Second * 5)
		_, err = GetCacheUint(key)
		if err == nil {
			t.Error("set cache want error but get nil\n")
		}
	})
}
func TestSetCacheBytes(t *testing.T) {
	defer Remove()
	defer zap_server.Remove()
	redisPwd := os.Getenv("redisPwd")
	CONFIG.Password = redisPwd
	t.Run("test set cache bytes", func(t *testing.T) {
		key := "test_set_cache"
		want := []byte("test_set_cache_value")
		err := SetCache(key, want, time.Duration(time.Second*3))
		if err != nil {
			t.Errorf("set cache get error %v\n", err)
		}
		get, err := GetCacheBytes(key)
		if err != nil {
			t.Errorf("set cache get error %v\n", err)
		}
		if !reflect.DeepEqual(get, want) {
			t.Errorf("set cache want [%s] but get [%s]\n", want, get)
		}
		time.Sleep(time.Second * 10)
		_, err = GetCacheBytes(key)
		if err == nil {
			t.Error("set cache want error but get nil\n")
		}
	})
}
