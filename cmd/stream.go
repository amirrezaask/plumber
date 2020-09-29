package main

import (
	"fmt"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/stream"
)

var natsStreamingRequiredKeys = [...]string{
	"url", "clusterID", "clientID", "subject",
}
var natsRequiredKeys = [...]string{}

func printerStreamFromConfig(c map[string]interface{}) (plumber.Stream, error) {
	return stream.NewPrinterStream()
}
func arrayStreamFromConfig(c map[string]interface{}) (plumber.Stream, error) {
	return stream.NewArrayStream(c["words"].([]interface{})...)
}
func natsStreamingFromConfig(c map[string]interface{}) (plumber.Stream, error) {
	for _, k := range natsStreamingRequiredKeys {
		v := c[k]
		if v == nil {
			return nil, fmt.Errorf("%s key is required for nats streaming.", k)
		}
	}
	s, err := stream.NewNatsStreaming(c["url"].(string), c["subject"].(string),
		c["clusterID"].(string), c["clientID"].(string))
	return s, err
}

func natsFromConfig(c map[string]interface{}) (plumber.Stream, error) {
	for _, k := range natsRequiredKeys {
		v := c[k]
		if v == nil {
			return nil, fmt.Errorf("%s key is required for nats.", k)
		}
	}
	s, err := stream.NewNats(c["url"].(string), c["subject"].(string))
	return s, err
}
