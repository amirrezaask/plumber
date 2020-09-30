package checkpoint

import (
	"time"

	"github.com/amirrezaask/plumber"
)

func WithInterval(d time.Duration) plumber.Checkpoint {
	return func(s plumber.Pipeline) {
		for range time.Tick(d) {
			s.UpdateState()
			s.State().Flush()
		}
	}
}
