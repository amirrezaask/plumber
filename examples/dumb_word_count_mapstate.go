package main

import (
	"fmt"
	"github.com/amirrezaask/plumber/pipe"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/checkpoint"
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
		c, err := ctx.State.Get(word)
		if err != nil {
			ctx.Err <- err
			return
		}
		if c == nil {
			c = 0
		}
		counter, ok := c.(int)
		if !ok {
			ctx.Logger.Errorf("state of %s is not an int", word)
			continue
		}
		counter = counter + 1
		err = ctx.State.Set(word, counter)
		if err != nil {
			ctx.Err <- err
			return
		}
		ctx.Out <- word
	}
}

func main() {
	s := state.NewMapState()

	errs, err := pipeline.
		NewDefaultSystem().
		WithLogger(logrus.New()).
		SetCheckpoint(checkpoint.WithInterval(time.Second * 1)).
		SetState(s).
		From(stream.NewArrayInput("amirreza", "parsa")).
		Then(toUpper).
		Then(pipe.MakePipe(toLower)).
		Then(count).
		To(stream.NewPrinterOutput()).
		Initiate()
	if err != nil {
		panic(err)
	}
	for err := range errs {
		fmt.Println(err)
	}
}
