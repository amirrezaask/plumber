package plumber

//Stream
type Stream interface {
	Write(interface{}) error
	StartReading() error
	ReadChan() chan interface{}
	// SinceLastCheckpoint() []string
}

// System
type System interface {
	Name() string
	State() State
	SetState(State) System
	Then(Lambda) System
	From(Stream) System
	To(Stream) System
	Initiate() (chan error, error)
}

//Each state backend should implement this.
type State interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

// Lambda is a stateful function
type Lambda func(state State, input interface{}) (interface{}, error)

//Checkpoints for fault tolerant system.
type Checkpoint func(State) error
