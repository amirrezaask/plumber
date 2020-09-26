package stream

import (
	"encoding/json"

	"github.com/amirrezaask/plumber"
	"github.com/nats-io/nats.go"
)

type NatsStream struct {
	nc       *nats.Conn
	subject  string
	readChan chan interface{}
}

func NewNatsStream(url string, subject string, options ...nats.Option) (plumber.Stream, error) {
	conn, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	return &NatsStream{
		nc:       conn,
		subject:  subject,
		readChan: make(chan interface{}),
	}, nil
}

func (n *NatsStream) ReadChan() chan interface{} {
	return n.readChan
}

func (n *NatsStream) StartReading() error {
	_, err := n.nc.Subscribe(n.subject, func(m *nats.Msg) {
		n.readChan <- string(m.Data)
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsStream) Write(v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return n.nc.Publish(n.subject, bs)
}
