package main

import (
	"fmt"
	"strings"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/state"
	"github.com/amirrezaask/plumber/stream"
	"github.com/amirrezaask/plumber/system"
)

func toLower(state plumber.State, value interface{}) (interface{}, error) {
	word := value.(string)
	word = strings.ToLower(word)
	return word, nil
}

func toUpper(state plumber.State, value interface{}) (interface{}, error) {
	word := value.(string)
	word = strings.ToUpper(word)
	return word, nil
}

func count(state plumber.State, input interface{}) (interface{}, error) {
	word := input.(string)
	counter, err := state.Get(string(word))
	if err != nil {
		return nil, err
	}
	if counter == nil {
		counter = 0
	}
	counter = counter.(int) + 1
	err = state.Set(string(word), counter)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func main() {
	state := state.NewMapState()
	input := stream.NewChanStream()
	output := stream.NewChanStream()
	system := system.NewDefaultSystem()
	system.SetState(state)
	errs, err := system.
		From(input).
		Then(toLower).
		Then(count).
		Then(toUpper).
		To(output).
		Initiate()
	if err != nil {
		panic("starting system failed")
	}
	for err := range errs {
		fmt.Println(err)
	}
}
