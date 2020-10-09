package pipe

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/amirrezaask/plumber"
)

//MakePipe create pipes from simple functions but using this will abstract away the data flow control.
func MakePipe(f func(state plumber.State, input interface{}) (interface{}, error)) plumber.Pipe {
	return func(ctx *plumber.PipeCtx) {
		for v := range ctx.In {
			o, err := f(ctx.State, v)
			if err != nil {
				ctx.Err <- err
				continue
			}
			ctx.Out <- o
		}
	}

}

//PipeFromExecutable creates a plumber.Pipe that uses given executable
func PipeFromExecutable(path string, needsState bool) plumber.Pipe {
	return func(ctx *plumber.PipeCtx) {
		args := []string{}
		if needsState {
			all, err := ctx.State.All()
			if err != nil {
				ctx.Err <- err
			}
			bs, err := json.Marshal(all)
			if err != nil {
				ctx.Err <- err
			}

			args = append(args, fmt.Sprintf("\"%s\"", string(bs)))
		}
		args = append(args, fmt.Sprintf("\"%v\"", <-ctx.In))

		c := exec.Command(path, args...)

		output, err := c.Output()
		if err != nil {
			ctx.Err <- err
		}
		//TODO: should update state cause based on the contract updated fileds are in output
		ctx.Out <- output
	}
}
