package dcache

import (
	"context"
	"errors"
	"time"

	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/memcache"
)

var (
	cacheDb = make(map[string][]byte)

	// ErrCacheMiss ...
	ErrCacheMiss = memcache.ErrCacheMiss
)

// Set ...
func Set(ctx context.Context, key string, value []byte) error {
	if appengine.IsAppEngine() {
		return memcache.Set(ctx, &memcache.Item{
			Key:   key,
			Value: value,
		})
	}
	cacheDb[key] = value
	return nil
}

// SetWithExpiration ...
func SetWithExpiration(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if appengine.IsAppEngine() {
		return memcache.Set(ctx, &memcache.Item{
			Key:        key,
			Value:      value,
			Expiration: expiration,
		})
	}
	cacheDb[key] = value
	return nil
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
	delete(cacheDb, key)
	return nil
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
	val, ok := cacheDb[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	return val, nil
}
