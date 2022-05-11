package internal

import (
	"bytes"
	"sync"
)

type safeBuffer struct {
	*bytes.Buffer
	sync.Mutex
}

func (sb *safeBuffer) Write(p []byte) (n int, err error) {
	sb.Lock()
	defer sb.Unlock()
	return sb.Buffer.Write(p)
}

func (sb *safeBuffer) String() string {
	sb.Lock()
	defer sb.Unlock()
	return sb.Buffer.String()
}

func (sb *safeBuffer) Reset() {
	sb.Lock()
	defer sb.Unlock()
	sb.Buffer.Reset()
}

func NewSafeBuffer() *safeBuffer {
	return &safeBuffer{
		Buffer: bytes.NewBuffer([]byte{}),
	}
}
