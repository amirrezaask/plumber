package stream

import (
	"fmt"
	"net/http"
)

type HttpStreamRequestToValue func(*http.Request) (interface{}, error)

type HttpStream struct {
	readChan       chan interface{}
	requestToValue HttpStreamRequestToValue
}

func NewHttpStream(requestToValue HttpStreamRequestToValue) (*HttpStream, error) {
	s := &HttpStream{
		readChan:       make(chan interface{}),
		requestToValue: requestToValue,
	}

	return s, nil
}

//http is stateless we don't need any state.
func (h *HttpStream) LoadState(map[string]interface{}) {
	return
}

func (h *HttpStream) Input() chan interface{} {
	return h.readChan
}

func (h *HttpStream) Output() chan interface{} {
	panic("Http stream can only be used as input")
}

//state is just empty
func (h *HttpStream) State() map[string]interface{} {
	return map[string]interface{}{}
}

func (h *HttpStream) Name() string {
	return "http-stream"
}

//Handler returns an HTTP handler func that you can use, pushing to stream happens in between before and after functions, you can pass either of them nil if you don't need anything to happen before or after.
func (h *HttpStream) Handler(before http.HandlerFunc, after http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if before != nil {
			before(w, r)
		}
		v, err := h.requestToValue(r)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Error on processing event: %v", err)
			return
		}
		h.readChan <- v
		if after != nil {
			after(w, r)
		}
	}
}
