package state

import (
	"sync"

	"github.com/amirrezaask/plumber"
)

type MapState struct {
	lock *sync.RWMutex
	m    map[string]interface{}
}

func (d *MapState) Get(k string) (interface{}, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	v := d.m[k]
	return v, nil
}

func (d *MapState) Set(k string, v interface{}) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.m[k] = v
	return nil
}
func (d *MapState) All() (map[string]interface{}, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	newM := make(map[string]interface{})
	for k, v := range d.m {
		if k == "____checkpoint" {
			continue
		}
		newM[k] = v
	}
	return newM, nil
}
func NewMapState() plumber.State {
	return &MapState{
		lock: &sync.RWMutex{},
		m:    make(map[string]interface{}),
	}
}
