package main

import (
	"context"

	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/state"
)

var redisRequiredKeys = [...]string{
	"host", "port",
}

func redisFromConfig(c map[string]interface{}) (plumber.State, error) {
	s, err := state.NewRedis(context.Background(), c["host"].(string), c["port"].(string), "", "", 0)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func mapFromConfig(c map[string]interface{}) (plumber.State, error) {
	return state.NewMapState(), nil
}
