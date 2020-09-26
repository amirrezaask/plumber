package stream

import "errors"

type DumbStream struct {
	Stream chan []byte
}

func (d *DumbStream) Read(p []byte) (n int, err error) {
	v := <-d.Stream
	if len(v) > len(p) {
		return -1, errors.New("size don't match")
	}
	for idx, b := range v {
		p[idx] = b
	}
	return len(p), nil
}
func (d *DumbStream) Write(p []byte) (n int, err error) {
	d.Stream <- p
	return len(p), nil
}
func (d *DumbStream) Close() error {
	close(d.Stream)
	return nil
}
