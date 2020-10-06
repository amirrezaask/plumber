package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/checkpoint"
	"github.com/amirrezaask/plumber/pipeline"
	"github.com/amirrezaask/plumber/state"
	"github.com/amirrezaask/plumber/stream"
)

func toLower(ctx *plumber.PipeCtx) {
	word := (<-ctx.In).(string)
	word = strings.ToLower(word)
	ctx.Out <- word
}

func toUpper(ctx *plumber.PipeCtx) {
	word := (<-ctx.In).(string)
	word = strings.ToUpper(word)
	ctx.Out <- word
}

func count(ctx *plumber.PipeCtx) {
	word := (<-ctx.In).(string)
	counter, err := ctx.State.GetInt(string(word))
	if err != nil {
		ctx.Err <- err
		return
	}
	counter = counter + 1
	err = ctx.State.Set(string(word), counter)
	if err != nil {
		ctx.Err <- err
		return
	}
	ctx.Out <- word
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
			input.Output() <- "salam"
		}
	}()
	output := stream.NewChanStream()
	//consume our output data
	go func() {
		for v := range output.Input() {
			if v != nil {
				fmt.Println(v)
			}
		}
	}()
	r, err := state.NewRedis(context.Background(), "localhost", "6379", "", "", 0)
	if err != nil {
		panic(err)
	}
	//create our plumber pipeline
	errs, err := pipeline.
		NewDefaultSystem().
		SetCheckpoint(checkpoint.WithInterval(time.Second * 1)).
		SetState(r).
		//SetState(state.NewBolt())
		// SetState(state.NewMapState()).
		From(input).
		Then(toLower).
		Then(toUpper).
		Then(count).
		To(output).
		Initiate()
	if err != nil {
		panic(err)
	}
	for err := range errs {
		fmt.Println(err)
	}
}
