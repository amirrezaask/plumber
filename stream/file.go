package stream

import (
	"encoding/json"
	"github.com/amirrezaask/plumber"
	"io"
	"os"
)

type fileInputState struct {
	CurrentByte uint64
}

type fileInput struct {
	fd           *os.File
	s            *fileInputState
	chunksLength uint64
	readChan     chan interface{}
}

func NewFileInput(path string, initialState *fileInputState, chunksLength uint64) (plumber.Input, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	f := &fileInput{
		fd:           fd,
		readChan:     make(chan interface{}),
		s:            initialState,
		chunksLength: chunksLength,
	}

	go func() {
		for {
			bs := make([]byte, f.chunksLength)
			n, err := f.fd.ReadAt(bs, int64(f.s.CurrentByte+1))
			if err != nil {
				return
			}
			f.s.CurrentByte += uint64(n)
			f.readChan <- bs
		}

	}()

	return f, nil
}
func (f *fileInput) Name() string {
	return "file-stream"
}

func (f *fileInput) LoadState(r io.Reader) error {
	s := &fileInputState{}
	return json.NewDecoder(r).Decode(s)
}

func (f *fileInput) Input() (chan interface{}, error) {
	return f.readChan, nil
}

func (f *fileInput) State() ([]byte, error) {
	return json.Marshal(f.s)
}

type fileOutputState struct {
	CurrentByte uint64
}

type fileOutput struct {
	fd           *os.File
	state        *fileOutputState
	chunksLength uint64
	writeChan    chan interface{}
}

func NewFileOutput(path string, initialState *fileOutputState, chunksLength uint64) (plumber.Output, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	f := &fileOutput{
		fd:           fd,
		state:        initialState,
		chunksLength: chunksLength,
	}
	go func() {
		for v := range f.writeChan {
			bs, err := json.Marshal(v)
			if err != nil {
				continue
			}
			n, err := f.fd.WriteAt(bs, int64(f.state.CurrentByte))
			if err != nil {
				continue
			}
			f.state.CurrentByte += uint64(n)
		}
	}()
	return f, nil
}

func (f *fileOutput) LoadState(reader io.Reader) error {
	s := &fileOutputState{}
	err := json.NewDecoder(reader).Decode(s)
	if err != nil {
		return err
	}
	f.state = s
	return nil
}

func (f *fileOutput) State() ([]byte, error) {
	return json.Marshal(f.state)
}

func (f *fileOutput) Name() string {
	return "file_output"
}

func (f *fileOutput) Output() (chan interface{}, error) {
	return f.writeChan, nil
}
