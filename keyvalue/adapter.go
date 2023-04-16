package keyvalue

// Adapter is a key-value adapter.
type Adapter interface {
	Get(string) (interface{}, bool)
	Put(string, interface{}) error
	Pairs(func(string, interface{}) bool) bool
}

// Copy copies the contents of `src` to `dst`.
func Copy(dst, src Adapter) error {
	var err error
	src.Pairs(func(key string, value interface{}) bool {
		err = dst.Put(key, value)
		return err == nil
	})
	return err
}
