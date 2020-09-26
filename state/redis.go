package state

import "github.com/amirrezaask/plumber"

type RedisState struct {
	conn
}

func NewRedisState(host, port, user, password, database string) plumber.State {

}
