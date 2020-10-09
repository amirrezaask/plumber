package pipe

import (
	"time"

	"github.com/amirrezaask/plumber"
)

func TimeWindow(d time.Duration) plumber.Pipe {
	return func(ctx *plumber.PipeCtx) {
		buff := []interface{}{}
		go func() {
			buff = append(buff, <-ctx.In)
		}()
		for range time.Tick(d) {
			ctx.Out <- buff
		}
	}
}
