package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/checkpoint"
	"github.com/amirrezaask/plumber/pipe"
	"github.com/amirrezaask/plumber/pipeline"
	"github.com/amirrezaask/plumber/state"
	"github.com/amirrezaask/plumber/stream"
)

func toLower(s plumber.State, i interface{}) (interface{}, error) {
	word := strings.ToLower(i.(string))
	return word, nil
}

func toUpper(ctx *plumber.PipeCtx) {
	for {
		word := (<-ctx.In).(string)
		word = strings.ToUpper(word)
		ctx.Out <- word
	}
}

func count(ctx *plumber.PipeCtx) {
	for {
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
}
func main() {

	r, err := state.NewRedis(context.Background(), "localhost", "6379", "", "", 0)
	if err != nil {
		panic(err)
	}
	//create our plumber pipeline
	errs, err := pipeline.
		NewDefaultSystem().
		SetCheckpoint(checkpoint.WithInterval(time.Second * 1)).
		SetState(r).
		From(stream.NewArrayStream("amirreza", "parsa")).
		Then(toUpper).
		Then(pipe.MakePipe(toLower)).
		Then(count).
		To(stream.NewPrinterStream()).
		Initiate()
	if err != nil {
		panic(err)
	}
	for err := range errs {
		fmt.Println(err)
	}
}
