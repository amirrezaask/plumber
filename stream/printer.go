package stream

import (
	"fmt"

	"github.com/amirrezaask/plumber"
)

type PrinterStream struct {
	writeChan chan interface{}
}

func NewPrinterStream() plumber.Stream {
	p := &PrinterStream{
		writeChan: make(chan interface{}),
	}
	go func() {
		for v := range p.writeChan {
			fmt.Printf("Printer:: %s\n", v)
		}
	}()
	return p
}
func (a *PrinterStream) LoadState(map[string]interface{}) {
	return
}
func (a *PrinterStream) Output() chan interface{} {
	return a.writeChan
}

func (a *PrinterStream) Input() chan interface{} {
	panic("printer should be only used as output of pipeline")
}

func (a *PrinterStream) State() map[string]interface{} {
	return nil
}

func (a *PrinterStream) Name() string {
	return "printer-stream"
}
