package plumber

type Opts map[string]interface{}

//Stream
type Stream interface {
	LoadState(map[string]interface{})
	Write(interface{}) error
	StartReading() error
	ReadChan() chan interface{}
	State() map[string]interface{}
	Name() string
}

//Pipeline
type Pipeline interface {
	UpdateState() error
	Errors() chan error
	SetCheckpoint(Checkpoint) Pipeline
	Checkpoint()
	Name() string
	State() State
	SetState(State) Pipeline
	Then(Pipe) Pipeline
	Thens(...Pipe) Pipeline
	From(Stream) Pipeline
	To(Stream) Pipeline
	Initiate() (chan error, error)
}

//Each state backend should implement this.
//TODO: state should have some kind of initial value registeration
type State interface {
	Set(key string, value interface{}) error
	GetInt(key string) (int, error)
	Get(key string) (interface{}, error)
	All() (map[string]interface{}, error)
	Flush() error
}

// Pipe is a stateful function
type Pipe func(state State, input interface{}) (interface{}, error)

//StreamConstructor is just a contract for all Streams to agree on.
type StreamConstrcutor func(opts map[string]interface{}) (Stream, error)

//Checkpoints for fault tolerant Pipeline.
type Checkpoint func(Pipeline)
