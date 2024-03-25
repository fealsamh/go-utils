package dcache

import (
	"context"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/memcache"
)

var (
	cacheDb *leveldb.DB

	// ErrCacheMiss ...
	ErrCacheMiss = memcache.ErrCacheMiss
)

func init() {
	if !appengine.IsAppEngine() {
		db, err := leveldb.OpenFile("localcache.db", nil)
		if err != nil {
			panic(err)
		}
		cacheDb = db
	}
}

// Set ...
func Set(ctx context.Context, key string, value []byte) error {
	if appengine.IsAppEngine() {
		return memcache.Set(ctx, &memcache.Item{
			Key:   key,
			Value: value,
		})
	}
	return cacheDb.Put([]byte(key), value, nil)
}

// Remove ...
func Remove(ctx context.Context, key string) error {
	if appengine.IsAppEngine() {
		if err := memcache.Delete(ctx, key); err != nil {
			if errors.Is(err, memcache.ErrCacheMiss) {
				return ErrCacheMiss
			}
			return err
		}
		return nil
	}
	return cacheDb.Delete([]byte(key), nil)
}

// Get ...
func Get(ctx context.Context, key string) ([]byte, error) {
	if appengine.IsAppEngine() {
		item, err := memcache.Get(ctx, key)
		if err != nil {
			if errors.Is(err, memcache.ErrCacheMiss) {
				return nil, ErrCacheMiss
			}
			return nil, err
		}
		return item.Value, nil
	}
	val, err := cacheDb.Get([]byte(key), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrCacheMiss
		}
		return nil, err
	}
	return val, nil
}
