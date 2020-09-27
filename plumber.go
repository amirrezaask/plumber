package plumber

//Stream
type Stream interface {
	Write(interface{}) error
	StartReading() error
	ReadChan() chan interface{}
	State() map[string]interface{}
}

// System
type System interface {
	Errors() chan error
	SetCheckpoint(Checkpoint) System
	Checkpoint()
	GetStateCopy() (map[string]interface{}, error)
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
	All() (map[string]interface{}, error)
}

// Lambda is a stateful function
type Lambda func(state State, input interface{}) (interface{}, error)

//Checkpoints for fault tolerant system.
type Checkpoint func(System)
