package stream

import (
	"encoding/json"
	"log"

	"github.com/amirrezaask/plumber"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type NatsStreaming struct {
	sc            stan.Conn
	readSubject   string
	writerSubject string
	writeChan     chan interface{}
	readChan      chan interface{}
	currentEvent  uint64
}

//TODO: update according to StreamConstructor
func NewNatsStreaming(url string, readSubject string, writeSubject string,
	clusterID string, clientID string, options ...stan.Option) (plumber.Stream, error) {
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
	n := &NatsStreaming{
		sc:            sc,
		readSubject:   readSubject,
		writerSubject: writeSubject,
		readChan:      make(chan interface{}),
		writeChan:     make(chan interface{}),
		currentEvent:  0,
	}
	//start reading
	n.sc.Subscribe(n.readSubject, func(msg *stan.Msg) {
		n.currentEvent++
		n.readChan <- msg.Data

	}, stan.DurableName("TODOCHANGE"))
	//start writing
	go func() {
		for v := range n.writeChan {
			bs, err := json.Marshal(v)
			if err != nil {
				log.Println(err)
			}
			err = n.sc.Publish(n.writerSubject, bs)
			if err != nil {
				log.Println(err)
			}

		}

	}()
	return n, nil
}

//Since we are using nats streaming and we use durable subscription we don't need any state.
func (n *NatsStreaming) State() map[string]interface{} {
	return map[string]interface{}{}
}
func (n *NatsStreaming) Name() string {
	return "Nats-Streaming"
}

func (n *NatsStreaming) LoadState(s map[string]interface{}) {
	n.currentEvent = s["current_event"].(uint64)
}
func (n *NatsStreaming) Input() chan interface{} {
	return n.readChan
}

func (n *NatsStreaming) Output() chan interface{} {
	return n.writeChan
}
