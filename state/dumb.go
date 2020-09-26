package state

import "github.com/amirrezaask/plumber"

type DumbState map[string]interface{}

func (d DumbState) Get(k string) (interface{}, error) {
	v := d[k]
	return v, nil
}

func (d DumbState) Set(k string, v interface{}) error {
	d[k] = v
	return nil
}

func NewDumbState() plumber.State {
	return DumbState{}
}
