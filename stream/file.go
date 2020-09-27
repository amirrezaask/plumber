package stream

import (
	"os"

	"github.com/amirrezaask/plumber"
)

type fileStream struct {
	fd           *os.File
	currentByte  uint64
	chunksLength uint64
	readChan     chan interface{}
}

func NewFileStream(path string, chunksLength uint64) (plumber.Stream, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &fileStream{
		fd:           fd,
		readChan:     make(chan interface{}),
		chunksLength: chunksLength,
	}, nil
}
func (f *fileStream) Write(v interface{}) error {
	_, err := f.fd.Write(v.([]byte))
	return err
}
func (f *fileStream) Name() string {
	return "file-stream"
}
func (f *fileStream) StartReading() error {
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
	return nil
}
func (f *fileStream) LoadState(s map[string]interface{}) {
	f.currentByte = s["current_byte"].(uint64)
}
func (f *fileStream) ReadChan() chan interface{} {
	return f.readChan
}

func (f *fileStream) State() map[string]interface{} {
	return map[string]interface{}{
		"current_byte": f.currentByte,
	}
}
