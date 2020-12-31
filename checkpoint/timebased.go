package checkpoint

import (
	"time"

	"github.com/amirrezaask/plumber"
)

func WithInterval(d time.Duration) plumber.Checkpoint {
	return func(s plumber.Pipeline) {
		for range time.Tick(d) {
			err := s.UpdateState()
			if err != nil {
				s.Logger().Errorf("checkpoint error: %s", err.Error())
				return
			}
			err = s.State().Flush()
			if err != nil {
				s.Logger().Errorf("checkpoint error: %s", err.Error())
			}
		}
	}
}
