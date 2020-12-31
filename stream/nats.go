package stream

import (
	"encoding/json"
	"github.com/amirrezaask/plumber"
	"github.com/nats-io/nats.go"
	"io"
)

type NatsInputState struct {
	CurrentEvent uint64
}

type NatsInput struct {
	sub          *nats.Subscription
	nc           *nats.Conn
	state        *NatsInputState
	currentEvent uint64
	readSubject  string
	readChan     chan interface{}
}

func NewNatsInput(initialState *NatsInputState, url string, readSubject string, options ...nats.Option) (plumber.Input, error) {
	conn, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	n := &NatsInput{
		nc:          conn,
		readSubject: readSubject,
		readChan:    make(chan interface{}),
	}
	sub, err := n.nc.Subscribe(readSubject, func(msg *nats.Msg) {
		n.currentEvent++
		n.readChan <- msg.Data
	})
	if err != nil {
		return nil, err
	}
	n.state = initialState
	n.sub = sub
	return n, nil
}
func (n *NatsInput) State() ([]byte, error) {
	return json.Marshal(n.state)
}

func (n *NatsInput) Input() (chan interface{}, error) {
	return n.readChan, nil
}
func (n *NatsInput) Name() string {
	return "nats-input"
}
func (n *NatsInput) LoadState(r io.Reader) error {
	s := &NatsInputState{}
	return json.NewDecoder(r).Decode(s)
}

type NatsOutputState struct {
	CurrentEvent uint64
}

type NatsOutput struct {
	nc           *nats.Conn
	state        *NatsOutputState
	currentEvent uint64
	writeSubject string
	writeChan    chan interface{}
}

func NewNatsOutput(initialState *NatsOutputState, url string, writeSubject string, options ...nats.Option) (plumber.Output, error) {
	conn, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	n := &NatsOutput{
		nc:           conn,
		writeSubject: writeSubject,
		writeChan:    make(chan interface{}),
	}
	n.state = initialState
	go func() {
		for v := range n.writeChan {
			//TODO: handle this error
			n.nc.Publish(writeSubject, v.([]byte))
		}
	}()
	return n, nil
}
func (n *NatsOutput) State() ([]byte, error) {
	return json.Marshal(n.state)
}

func (n *NatsOutput) Output() (chan interface{}, error) {
	return n.writeChan, nil
}
func (n *NatsOutput) Name() string {
	return "nats-output"
}
func (n *NatsOutput) LoadState(r io.Reader) error {
	s := &NatsOutputState{}
	return json.NewDecoder(r).Decode(s)
}
