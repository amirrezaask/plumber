package main

import (
	"fmt"
	"strings"

	"github.com/amirrezaask/plumber"
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
	state := &plumber.DumbState{}
	source := make(chan interface{})
	go func() {
		for {
			source <- "Hello ThIIs iS plumber"
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
	errs := system.From(source).Map(toLower).Map(count).Map(toUpper).To(sink).Initiate()
	for err := range errs {
		fmt.Println(err)
	}
}
