package stream

import (
	"fmt"
	"io"

	"github.com/amirrezaask/plumber"
)

type PrinterOutput struct {
	writeChan chan interface{}
}

func NewPrinterOutput() plumber.Output {
	p := &PrinterOutput{
		writeChan: make(chan interface{}),
	}
	go func() {
		for v := range p.writeChan {
			fmt.Printf("Printer:: %s\n", v)
		}
	}()
	return p
}

func (a *PrinterOutput) LoadState(r io.Reader) error {
	return nil
}

func (a *PrinterOutput) Output() (chan interface{}, error) {
	return a.writeChan, nil
}

func (a *PrinterOutput) State() ([]byte, error) {
	return nil, nil
}

func (a *PrinterOutput) Name() string {
	return "printer-stream"
}
