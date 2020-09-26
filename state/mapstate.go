package state

import "github.com/amirrezaask/plumber"

type MapState map[string]interface{}

func (d MapState) Get(k string) (interface{}, error) {
	v := d[k]
	return v, nil
}

func (d MapState) Set(k string, v interface{}) error {
	d[k] = v
	return nil
}
func (d MapState) All() (map[string]interface{}, error) {
	return d, nil
}
func NewMapState() plumber.State {
	return MapState{}
}
