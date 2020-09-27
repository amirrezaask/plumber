package checkpoint

import (
	"fmt"
	"time"

	"github.com/amirrezaask/plumber"
)

func WithInterval(d time.Duration) plumber.Checkpoint {
	return func(s plumber.System) {
		for range time.Tick(d) {
			m, err := s.GetStateCopy()
			if err != nil {
				s.Errors() <- err
				continue
			}
			s.State().Set("____checkpoint", m)
			fmt.Printf("\n\n%+v", m)
		}
	}
}
