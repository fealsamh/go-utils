package keyvalue

// Adapter is a key-value adapter.
type Adapter interface {
	Get(string) (interface{}, bool)
	Put(string, interface{}) bool
	Pairs(func(string, interface{}) bool) bool
}

// Copy copies the contents of `src` to `dst`.
func Copy(dst, src Adapter) bool {
	return src.Pairs(func(key string, value interface{}) bool {
		return dst.Put(key, value)
	})
}
