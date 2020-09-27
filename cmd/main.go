package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/checkpoint"
	"github.com/amirrezaask/plumber/system"
)

func lambdaFromBin(path string) plumber.Lambda {
	return func(s plumber.State, input interface{}) (interface{}, error) {
		all, err := s.All()
		if err != nil {
			return nil, err
		}
		bs, err := json.Marshal(all)
		if err != nil {
			return nil, err
		}
		c := exec.Command(path,
			fmt.Sprintf("\"%s\"", string(bs)),
			fmt.Sprintf("\"%v\"", input))
		output, err := c.Output()
		if err != nil {
			return nil, err
		}
		//TODO: should update state cause based on the contract updated fileds are in output
		return string(output), nil
	}
}

type Stream struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}
type Checkpoint struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}
type State struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

type config struct {
	From       *Stream     `json:"from"`
	To         *Stream     `json:"to"`
	Checkpoint *Checkpoint `json:"checkpoint"`
	State      *State      `json:"state"`
	Pipeline   []string    `json:"pipeline"`
}

func (c *config) toStream() (plumber.Stream, error) {
	switch c.To.Type {
	case "nats":
		return natsFromConfig(c.From.Args)
	case "nats-streaming":
		return natsStreamingFromConfig(c.From.Args)
	default:
		return nil, errors.New("not found")
	}
}
func (c *config) fromStream() (plumber.Stream, error) {
	switch c.From.Type {
	case "nats":
		return natsFromConfig(c.From.Args)
	case "nats-streaming":
		return natsStreamingFromConfig(c.From.Args)
	default:
		return nil, errors.New("not found")
	}
}

func (c *config) state() (plumber.State, error) {
	switch c.State.Type {
	case "redis":
		return redisFromConfig(c.State.Args)
	default:
		return nil, errors.New("not found")
	}

}

func (c *config) checkpoint() (plumber.Checkpoint, error) {
	switch c.Checkpoint.Type {
	case "time-based":
		return checkpoint.WithInterval(time.Second * 2), nil
	default:
		return nil, errors.New("notfound")
	}
}
func main() {
	if len(os.Args) < 2 {
		log.Fatal("need config")
	}
	bs, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var c config
	err = json.Unmarshal(bs, &c)
	if err != nil {
		log.Fatal(err)
	}
	s := system.NewDefaultSystem()
	// configure from stream
	from, err := c.fromStream()
	if err != nil {
		log.Fatal(err)
	}
	s.From(from)
	// configure to stream
	to, err := c.toStream()
	if err != nil {
		log.Fatal(err)
	}
	s.To(to)
	// Set state backend
	st, err := c.state()
	if err != nil {
		log.Fatal(err)
	}
	s.SetState(st)
	// Set checkpoints handler
	chpt, err := c.checkpoint()
	if err != nil {
		log.Fatal(err)
	}
	s.SetCheckpoint(chpt)
	// configure pipes
	pipes := []plumber.Lambda{}
	for _, p := range c.Pipeline {
		pipes = append(pipes, lambdaFromBin(p))
	}
	s.Thens(pipes...)

	errs, err := s.Initiate()
	if err != nil {
		log.Fatal(err)
	}

	for err := range errs {
		fmt.Println(err)
	}

}
