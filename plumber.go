package plumber

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

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

//PipeFromBin creats a Pipe from given pipe object.
func PipeFromExecutable(path string, needsState bool) Pipe {
	return func(s State, input interface{}) (interface{}, error) {
		all, err := s.All()
		if err != nil {
			return nil, err
		}
		bs, err := json.Marshal(all)
		if err != nil {
			return nil, err
		}

		args := []string{}
		if needsState {
			args = append(args, fmt.Sprintf("\"%s\"", string(bs)))
		}
		args = append(args, fmt.Sprintf("\"%v\"", input))

		c := exec.Command(path, args...)

		output, err := c.Output()
		if err != nil {
			return nil, err
		}
		//TODO: should update state cause based on the contract updated fileds are in output
		return string(output), nil
	}
}

//Checkpoints for fault tolerant Pipeline.
type Checkpoint func(Pipeline)
