package dcache

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"
)

// SetObject ...
func SetObject[T any](ctx context.Context, key string, obj *T) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(obj); err != nil {
		return err
	}
	return Set(ctx, key, buf.Bytes())
}

// SetObjectWithExpiration ...
func SetObjectWithExpiration[T any](ctx context.Context, key string, obj *T, expiration time.Duration) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(obj); err != nil {
		return err
	}
	return SetWithExpiration(ctx, key, buf.Bytes(), expiration)
}

// GetObject ...
func GetObject[T any](ctx context.Context, key string) (*T, error) {
	b, err := Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var obj T
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
