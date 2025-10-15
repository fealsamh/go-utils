package kvo

// Coding is the KVC interface.
type Coding interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
}

// ObservingTime determines the observing time point.
type ObservingTime int

const (
	// Prior observes before setting a new value.
	Prior ObservingTime = iota
	// Post observes after setting a new value.
	Post
)

// Observed represents an observable object.
type Observed interface {
	AddObserver(key string, ot ObservingTime, observer func(interface{}, interface{}))
}

// Object is a dynamic observable object.
type Object struct {
	values         map[string]interface{}
	observersPrior []func(interface{}, interface{})
	observersPost  []func(interface{}, interface{})
}

// NewObject creates a new [Object].
func NewObject() *Object {
	return &Object{values: make(map[string]interface{})}
}

// Set sets a property.
func (obj *Object) Set(key string, value interface{}) bool {
	oldVal := obj.values[key]
	for _, obs := range obj.observersPrior {
		obs(oldVal, value)
	}
	obj.values[key] = value
	for _, obs := range obj.observersPost {
		obs(oldVal, value)
	}
	return true
}

// Get gets a property.
func (obj *Object) Get(key string) (interface{}, bool) {
	val, ok := obj.values[key]
	return val, ok
}

// AddObserver adds an observer.
func (obj *Object) AddObserver(key string, ot ObservingTime, observer func(interface{}, interface{})) {
	switch ot {
	case Prior:
		obj.observersPrior = append(obj.observersPrior, observer)
	case Post:
		obj.observersPost = append(obj.observersPost, observer)
	}
}
