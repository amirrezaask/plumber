package plumber

type Stream interface {
	LoadState(map[string]interface{}) error
	State() map[string]interface{}
	Name() string
}

type Input interface {
	Stream
	Input() (chan interface{}, error)
}

type Output interface {
	Stream
	Output() (chan interface{}, error)
}

//Pipeline handles data flow between streams also handling checkpoints and state.
type Pipeline interface {
	Logger() Logger
	UpdateState() error
	Errors() chan error
	SetCheckpoint(Checkpoint) Pipeline
	Checkpoint()
	Name() string
	State() State
	SetState(State) Pipeline
	Then(Pipe) Pipeline
	Thens(...Pipe) Pipeline
	From(Input) Pipeline
	To(Output) Pipeline
	Initiate() (chan error, error)
	GetInputStream() Input
	GetOutputStream() Output
	WithLogger(l Logger)
}

//State
type State interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	All() (map[string]interface{}, error)
	Flush() error
}

//PipeCtx is the only way a pipe can talk to outside world.
type PipeCtx struct {
	State  State
	In     chan interface{}
	Out    chan interface{}
	Err    chan error
	Logger Logger
}

// Pipe is a stateful function
type Pipe func(*PipeCtx)

//Checkpoints for fault tolerant Pipeline.
type Checkpoint func(Pipeline)
