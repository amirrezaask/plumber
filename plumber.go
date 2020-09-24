package plumber

type Stream chan interface{}

//Each state backend should implement this.
type State interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

// Lambda is a stateful function
type Lambda func(state State, input interface{}) (interface{}, error)

type lamdaContainer struct {
	l   Lambda
	In  Stream
	Out Stream
}
