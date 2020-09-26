package stream

import (
	"encoding/json"

	"github.com/amirrezaask/plumber"
	"github.com/labstack/gommon/log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type NatsStreaming struct {
	sc       stan.Conn
	subject  string
	readChan chan interface{}
}

func NewNatsStreaming(url string, subject string,
	clusterID string, clientID string, options ...stan.Option) (plumber.Stream, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Errorf("Connection lost, reason: %v", reason)
		}))

	if err != nil {
		return nil, err
	}
	return &NatsStreaming{
		sc:       sc,
		subject:  subject,
		readChan: make(chan interface{}),
	}, nil
}

//Since we are using nats streaming and we use durable subscription we don't need any state.
func (n *NatsStreaming) State() map[string]interface{} {
	return map[string]interface{}{}
}

func (n *NatsStreaming) ReadChan() chan interface{} {
	return n.readChan
}

func (n *NatsStreaming) StartReading() error {
	_, err := n.sc.Subscribe(n.subject, func(msg *stan.Msg) {
		n.ReadChan() <- string(msg.Data)
	}, stan.DurableName("PLUMBER_DURABLE")) //TODO: variable for durable name
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsStreaming) Write(v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return n.sc.Publish(n.subject, bs)
}
