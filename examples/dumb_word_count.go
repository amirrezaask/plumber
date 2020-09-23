package main

import (
	"fmt"

	"github.com/amirrezaask/plumber"
)

func lambda(state plumber.State, in plumber.Stream, _ plumber.Stream, errs plumber.Stream) {
start:
	word := make([]byte, 1024)
	_, err := in.Read(word)
	if err != nil {
		errs.Write([]byte(err.Error()))
		goto start
	}
	counter, err := state.Get(string(word))
	if err != nil {
		errs.Write([]byte(err.Error()))
		goto start
	}
	if counter == nil {
		counter = 0
	}
	counter = counter.(int) + 1
	err = state.Set(string(word), counter)
	if err != nil {
		errs.Write([]byte(err.Error()))
		goto start
	}
	goto start
}
func main() {
	state := plumber.DumbState{}
	in := &plumber.DumbStream{
		Stream: make(chan []byte),
	}
	go func() {
		for {
			_, err := in.Write([]byte("This is plumber"))
			if err != nil {
				fmt.Println("err  ", err)
			}
		}
	}()
	errs := &plumber.DumbStream{
		Stream: make(chan []byte),
	}
	go func() {
		var err []byte
		for {
			_, errRead := errs.Read(err)
			if err != nil {
				fmt.Println("err:: ", errRead)
			}
			if len(err) > 0 {
				fmt.Println("err >> ", err)
			}
		}
	}()
	lambda(state, in, nil, errs)
}
