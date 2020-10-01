package stream

import (
	"fmt"
	"net/http"

	"github.com/amirrezaask/plumber"
)

type HttpStreamRequestToValue func(*http.Request) (interface{}, error)

type httpStream struct {
	readChan       chan interface{}
	requestToValue HttpStreamRequestToValue
}

func NewHttpStream(requestToValue HttpStreamRequestToValue) (plumber.Stream, error) {
	return &httpStream{
		readChan:       make(chan interface{}),
		requestToValue: requestToValue,
	}, nil
}

//http is stateless we don't need any state.
func (h *httpStream) LoadState(map[string]interface{}) {
	return
}

//Write does not do anything since it's http after all.
func (h *httpStream) Write(interface{}) error {
	return nil
}

//Since http is our input we don't need start a goroutine to read from it.
func (h *httpStream) StartReading() error {
	return nil
}

func (h *httpStream) ReadChan() chan interface{} {
	return h.readChan
}

//state is just empty
func (h *httpStream) State() map[string]interface{} {
	return map[string]interface{}{}
}

func (h *httpStream) Name() string {
	return "http-stream"
}

//Handler returns an HTTP handler func that you can use, pushing to stream happens in between before and after functions, you can pass either of them nil if you don't need anything to happen before or after.
func (h *httpStream) Handler(before http.HandlerFunc, after http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
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
