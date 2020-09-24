package main

import (
	"fmt"

	"github.com/amirrezaask/plumber"
)

func lambda(state plumber.State, value interface{}) (interface{}, error) {
	word := value.(string)
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
	return nil, nil
}
func main() {
	state := &plumber.DumbState{}
	source := make(chan interface{})
	go func() {
		for {
			source <- "Hello This is plumber"
		}
	}()
	sink := make(chan interface{})
	go func() {
		for word := range sink {
			fmt.Println(word)
		}
	}()
	system := plumber.NewDefaultSystem()
	system.SetState(state)
	errs := system.From(source).Map(lambda).To(sink).Initiate()
	for err := range errs {
		fmt.Println(err)
	}
}
