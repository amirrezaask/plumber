package stream

import (
	"encoding/json"
	"log"

	"github.com/amirrezaask/plumber"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	nc           *nats.Conn
	currentEvent uint64
	readSubject  string
	writeSubject string
	readChan     chan interface{}
	writeChan    chan interface{}
}

//TODO: update according to stream constructor.
func NewNats(url string, readSubject string, writeSubject string, options ...nats.Option) (plumber.Stream, error) {
	conn, err := nats.Connect(url, options...)
	if err != nil {
		return nil, err
	}
	n := &Nats{
		nc:           conn,
		readSubject:  readSubject,
		writeSubject: writeSubject,
		readChan:     make(chan interface{}),
		writeChan:    make(chan interface{}),
	}
	n.nc.Subscribe(readSubject, func(msg *nats.Msg) {
		n.currentEvent++
		n.readChan <- msg.Data
	})
	go func() {
		for v := range n.writeChan {
			bs, err := json.Marshal(v)
			if err != nil {
				log.Println(err)
				continue
			}
			err = n.nc.Publish(n.writeSubject, bs)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}()
	return n, nil
}
func (n *Nats) State() map[string]interface{} {
	return map[string]interface{}{}
}

func (n *Nats) Output() chan interface{} {
	return n.readChan
}
func (n *Nats) Input() chan interface{} {
	return n.readChan
}
func (n *Nats) Name() string {
	return "NATS"
}
func (n *Nats) LoadState(s map[string]interface{}) {}
