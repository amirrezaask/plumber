package stream

import (
	"encoding/json"

	"github.com/amirrezaask/plumber"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	nc       *nats.Conn
	subject  string
	readChan chan interface{}
}

func NewNats(url string, subject string, options ...nats.Option) (plumber.Stream, error) {
	conn, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	return &Nats{
		nc:       conn,
		subject:  subject,
		readChan: make(chan interface{}),
	}, nil
}
func (n *Nats) State() map[string]interface{} {
	return map[string]interface{}{}
}

func (n *Nats) ReadChan() chan interface{} {
	return n.readChan
}

func (n *Nats) StartReading() error {
	_, err := n.nc.Subscribe(n.subject, func(m *nats.Msg) {
		n.readChan <- string(m.Data)
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *Nats) Write(v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return n.nc.Publish(n.subject, bs)
}
