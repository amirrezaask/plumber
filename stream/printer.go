package stream

import (
	"fmt"

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

func (a *PrinterOutput) LoadState(map[string]interface{}) error {
	return nil
}

func (a *PrinterOutput) Output() (chan interface{}, error) {
	return a.writeChan, nil
}

func (a *PrinterOutput) State() map[string]interface{} {
	return nil
}

func (a *PrinterOutput) Name() string {
	return "printer-stream"
}
