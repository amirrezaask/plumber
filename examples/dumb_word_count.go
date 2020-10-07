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
		Then(toLower).
		Then(toUpper).
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
