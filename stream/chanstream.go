package stream

import "github.com/amirrezaask/plumber"

type ChanStream struct {
	c         chan interface{}
	readChan  chan interface{}
	writeChan chan interface{}
}

//TODO: update according to stream constructor
func NewChanStream() plumber.Stream {
	st := &ChanStream{c: make(chan interface{}), readChan: make(chan interface{})}
	go func() {
		for v := range st.c {
			st.Input() <- v
		}
	}()

	go func() {
		for v := range st.writeChan {
			st.c <- v
		}
	}()
	return st
}

func (d *ChanStream) Name() string {
	return "chan-stream"
}

func (d *ChanStream) LoadState(m map[string]interface{}) {
	return
}
func (d *ChanStream) State() map[string]interface{} {
	return map[string]interface{}{}
}
func (d *ChanStream) Input() chan interface{} {
	return d.readChan
}

func (d *ChanStream) Output() chan interface{} {
	return d.writeChan
}
