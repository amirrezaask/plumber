package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/amirrezaask/plumber"
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
	fmt.Printf("%+v", c)
}
