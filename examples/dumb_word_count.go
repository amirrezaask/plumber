package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/checkpoint"
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
	// input, err := stream.NewNatsStreaming("localhost:4222", "plumber", "clusterID", "thisclient")
	// input, err := stream.NewNats("localhost:4222", "plumber")
	// if err != nil {
	// 	panic(err)
	// }
	input := stream.NewChanStream()
	// feed some data into our input stream
	go func() {
		for {
			input.Write("This Is tHe eNd")
		}
	}()
	output := stream.NewChanStream()
	//consume our output data
	go func() {
		for v := range output.ReadChan() {
			if v != nil {
				fmt.Println(v)
			}
		}
	}()

	//create our plumber pipeline
	errs, err := system.
		NewDefaultSystem().
		SetCheckpoint(checkpoint.WithInterval(time.Second * 1)).
		//SetState(state.NewRedis())
		//SetState(state.NewBolt())
		SetState(state.NewMapState()).
		From(input).
		Then(toLower).
		Then(toUpper).
		Then(count).
		To(output).
		Initiate()
	if err != nil {
		panic("starting system failed")
	}
	for err := range errs {
		fmt.Println(err)
	}
}
