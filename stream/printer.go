package stream

import (
	"fmt"

	"github.com/amirrezaask/plumber"
)

type PrinterStream struct {
}

func NewPrinterStream() (plumber.Stream, error) {
	return &PrinterStream{}, nil
}
func (a *PrinterStream) LoadState(map[string]interface{}) {
	return
}
func (a *PrinterStream) Write(v interface{}) error {
	fmt.Printf("PrinterStream => %v", v)
	return nil
}
func (a *PrinterStream) StartReading() error {
	return nil
}
func (a *PrinterStream) ReadChan() chan interface{} {
	return nil
}

func (a *PrinterStream) State() map[string]interface{} {
	return nil
}

func (a *PrinterStream) Name() string {
	return "printer-stream"
}
