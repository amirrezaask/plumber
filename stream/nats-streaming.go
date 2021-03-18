package stream

import (
	"encoding/json"
	"io"
	"log"

	"github.com/amirrezaask/plumber"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type NatsStreamingInputState struct {
	DurableName string
}

type NatsStreamingInput struct {
	sc           stan.Conn
	sub          stan.Subscription
	readSubject  string
	readChan     chan interface{}
	currentEvent uint64
}

func NewNatsStreamingInput(initialState *NatsStreamingInputState, url string, readSubject string, clusterID string, clientID string, options ...stan.Option) (plumber.Input, error) {
	//nc, err := nats.Connect(url)
	//if err != nil {
	//	return nil, err
	//}
	sc, err := stan.Connect(clusterID, clientID)

	if err != nil {
		return nil, err
	}
	n := &NatsStreamingInput{
		sc:           sc,
		readSubject:  readSubject,
		readChan:     make(chan interface{}),
		currentEvent: 0,
	}
	//start reading
	sub, err := n.sc.Subscribe(n.readSubject, func(msg *stan.Msg) {
		n.currentEvent++
		n.readChan <- msg.Data

	}, stan.DurableName(initialState.DurableName))
	n.sub = sub
	return n, nil
}

//Since we are using nats streaming and we use durable subscription we don't need any state.
func (n *NatsStreamingInput) State() ([]byte, error) {
	return nil, nil
}
func (n *NatsStreamingInput) Name() string {
	return "NatsStreaming-Input"
}

//Since we are using nats streaming and we use durable subscription we don't need any state.
func (n *NatsStreamingInput) LoadState(r io.Reader) error {
	state := &NatsStreamingInputState{}
	return json.NewDecoder(r).Decode(state)
}

func (n *NatsStreamingInput) Input() (chan interface{}, error) {
	return n.readChan, nil
}

type NatsStreamingOutputState struct {
	DurableName string
}

type NatsStreamingOutput struct {
	sc           stan.Conn
	writeSubject string
	writeChan    chan interface{}
	currentEvent uint64
}

func NewNatsStreamingOutput(initialState *NatsStreamingOutputState, url string, writeSubject string, clusterID string, clientID string, options ...stan.Option) (plumber.Output, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Println("Connection lost, reason: %v", reason)
			return
		}))

	if err != nil {
		return nil, err
	}
	n := &NatsStreamingOutput{
		sc:           sc,
		writeSubject: writeSubject,
		writeChan:    make(chan interface{}),
		currentEvent: 0,
	}
	go func() {
		for v := range n.writeChan {
			//TODO: handle this error
			n.sc.Publish(writeSubject, v.([]byte))
		}
	}()
	return n, nil
}

//Since we are using nats streaming and we use durable subscription we don't need any state.
func (n *NatsStreamingOutput) State() ([]byte, error) {
	return nil, nil
}
func (n *NatsStreamingOutput) Name() string {
	return "NatsStreaming-Output"
}

//Since we are using nats streaming and we use durable subscription we don't need any state.
func (n *NatsStreamingOutput) LoadState(r io.Reader) error {
	state := &NatsStreamingInputState{}
	return json.NewDecoder(r).Decode(state)
}

func (n *NatsStreamingOutput) Output() (chan interface{}, error) {
	return n.writeChan, nil
}
