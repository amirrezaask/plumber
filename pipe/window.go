package pipe

import (
	"time"

	"github.com/amirrezaask/plumber"
)

func TimeWindow(d time.Duration) plumber.Pipe {
	return func(ctx *plumber.PipeCtx) {
		buff := []interface{}{}
		go func() {
			for v := range ctx.In {
				buff = append(buff, v)
			}
		}()
		for range time.Tick(d) {
			ctx.Out <- buff
		}
	}
}

func CountWindow(size int) plumber.Pipe {
	return func(ctx *plumber.PipeCtx) {
		buff := []interface{}{}
		release := make(chan struct{})
		go func() {
			for v := range ctx.In {
				buff = append(buff, v)
				if len(buff) == size {
					release <- struct{}{}
				}
			}
		}()
		for range release {
			ctx.Out <- buff
			//TODO: check if a data race happens
			buff = []interface{}{}
		}
	}

}

func SignalWindow(sig chan interface{}) plumber.Pipe {
	return func(ctx *plumber.PipeCtx) {
		buff := []interface{}{}
		go func() {
			for v := range ctx.In {
				buff = append(buff, v)
			}
		}()
		for range sig {
			ctx.Out <- buff
			//TODO: check if a data race happens
			buff = []interface{}{}
		}
	}
}
