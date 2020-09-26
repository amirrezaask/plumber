package stream

import "github.com/amirrezaask/plumber"

func NewDumbStream(handler func(plumber.Stream)) plumber.Stream {
	c := make(chan interface{})
	go handler(c)
	return c
}
