package checkpoint

import (
	"time"

	"github.com/amirrezaask/plumber"
)

func WithInterval(d time.Duration) plumber.Checkpoint {
	return func(s plumber.State) error {
		return s.Set("checkpoint", s)
	}
}
