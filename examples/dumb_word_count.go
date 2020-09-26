package main

import (
	"fmt"
	"strings"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/state"
	"github.com/amirrezaask/plumber/stream"
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
	state := state.NewDumbState()
	input := stream.NewDumbStream(func(s plumber.Stream) {
		for {
			s <- "Hello ThIIs iS plumber"
		}
	})
	output := stream.NewDumbStream(func(s plumber.Stream) {
		for v := range s {
			if v != nil {
				fmt.Println(v)
			}
		}
	})
	system := plumber.NewDefaultSystem()
	system.SetState(state)
	errs := system.
		Eat(input).
		Digest(toLower).
		Digest(count).
		Digest(toUpper).
		Defecate(output).
		Initiate()
	for err := range errs {
		fmt.Println(err)
	}
}
