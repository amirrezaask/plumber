package checkpoint

import (
	"fmt"
	"time"

	"github.com/amirrezaask/plumber"
)

func WithInterval(d time.Duration) plumber.Checkpoint {
	return func(s plumber.State) error {
		for range time.Tick(d) {
			fmt.Println("$$$$$$$$$$")
			all, _ := s.All()
			fmt.Printf("%v\n", all)
			s.Set("____checkpoint", all)
			fmt.Println("###########")
		}
		return nil
	}
}
