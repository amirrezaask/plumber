package stream

import (
	"encoding/json"
	"os"

	"github.com/amirrezaask/plumber"
)

type fileStream struct {
	fd           *os.File
	currentByte  uint64
	chunksLength uint64
	readChan     chan interface{}
	writeChan    chan interface{}
}

func NewFileStream(path string, chunksLength uint64) (plumber.Stream, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	f := &fileStream{
		fd:           fd,
		readChan:     make(chan interface{}),
		chunksLength: chunksLength,
		writeChan:    make(chan interface{}),
		currentByte:  0,
	}

	go func() {
		for {
			bs := make([]byte, f.chunksLength)
			n, err := f.fd.ReadAt(bs, int64(f.currentByte+1))
			if err != nil {
				return
			}
			f.currentByte += uint64(n)
			f.readChan <- bs
		}

	}()

	go func() {
		for v := range f.writeChan {
			bs, err := json.Marshal(v)
			if err != nil {
				continue
			}
			n, err := f.fd.WriteAt(bs, int64(f.currentByte))
			if err != nil {
				continue
			}
			f.currentByte += uint64(n)
		}
	}()

	return f, nil
}
func (f *fileStream) Name() string {
	return "file-stream"
}

func (f *fileStream) LoadState(s map[string]interface{}) {
	f.currentByte = s["current_byte"].(uint64)
}

func (f *fileStream) Input() chan interface{} {
	return f.readChan
}

func (f *fileStream) Output() chan interface{} {
	return f.writeChan
}

func (f *fileStream) State() map[string]interface{} {
	return map[string]interface{}{
		"current_byte": f.currentByte,
	}
}
