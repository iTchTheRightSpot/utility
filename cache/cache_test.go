package cache

import (
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/utility/utils"
	"reflect"
	"sync"
	"testing"
	"time"
)

type cxObj struct {
	name string
}

func TestInMemoryCache(t *testing.T) {
	t.Parallel()

	t.Run("should insert and retrieve", func(t *testing.T) {
		t.Parallel()

		cache := InMemoryCache[string, cxObj]{
			logger:   utils.DevLogger("UTC"),
			cache:    &sync.Map{},
			duration: time.Duration(2) * time.Second,
			size:     2,
		}

		// given
		key := uuid.NewString()
		obj := cxObj{name: "hello world"}

		// methods to test
		cache.Put(key, obj)
		val := cache.Get(key)

		// assert
		if val == nil {
			t.Errorf("expect %s given nil", val)
		}

		if !reflect.DeepEqual(obj, *val) {
			t.Errorf("expect %s to equal given %v", obj, *val)
		}
	})

	t.Run("should insert and self delete", func(t *testing.T) {
		t.Parallel()

		cache := InMemoryCache[string, cxObj]{
			logger:   utils.DevLogger("UTC"),
			cache:    &sync.Map{},
			duration: time.Duration(2) * time.Second,
			size:     2,
		}

		// given & method to test
		key := uuid.NewString()
		cache.Put(key, cxObj{name: "hello world"})
		time.Sleep(time.Duration(3) * time.Second)

		// assert
		val := cache.Get(key)
		if val != nil {
			t.Errorf("expect %s given nil", val)
		}
	})

	t.Run("should insert & delete", func(t *testing.T) {
		t.Parallel()

		cache := InMemoryCache[string, cxObj]{
			logger:   utils.DevLogger("UTC"),
			cache:    &sync.Map{},
			duration: time.Duration(2) * time.Second,
			size:     2,
		}

		// given
		key := uuid.NewString()

		// methods to test
		cache.Put(key, cxObj{name: "custom object"})
		cache.Delete(key)
		val := cache.Get(key)

		// assert
		if val != nil {
			t.Errorf("expect nil given %s", val)
		}
	})

	t.Run("validate max size not exceeded", func(t *testing.T) {
		t.Parallel()

		cache := InMemoryCache[string, cxObj]{
			logger:   utils.DevLogger("UTC"),
			cache:    &sync.Map{},
			duration: time.Duration(2) * time.Second,
			size:     2,
		}

		// methods to test
		cache.Put(uuid.NewString(), cxObj{name: "hello world 1"})
		cache.Put(uuid.NewString(), cxObj{name: "hello world 2"})
		cache.Put(uuid.NewString(), cxObj{name: "hello world 3"})

		size := cache.Length()
		if size > cache.size {
			t.Errorf("expect %v given %v", cache.size, size)
		}

		if size != cache.size {
			t.Errorf("expect %v given %v", cache.size, size)
		}
	})

	t.Run("clear all", func(t *testing.T) {
		t.Parallel()

		cache := InMemoryCache[string, cxObj]{
			logger:   utils.DevLogger("UTC"),
			cache:    &sync.Map{},
			duration: time.Duration(2) * time.Second,
			size:     2,
		}

		// given
		cache.Put(uuid.NewString(), cxObj{name: "hello world 1"})
		cache.Put(uuid.NewString(), cxObj{name: "hello world 2"})

		// method to test
		cache.Clear()

		// assert
		size := cache.Length()
		if size != 0 {
			t.Errorf("expect 0 given %v", size)
		}
	})
}
