package main

import (
	"bytes"
	"io"
	"sync"
)

// prefixed wraps an io.Writer so each line gets a "[name] " prefix. Buffers
// partial lines so a tag never lands mid-line.
func prefixed(name string, w io.Writer) io.Writer {
	return &lineTagger{name: name, out: w}
}

type lineTagger struct {
	name string
	out  io.Writer
	mu   sync.Mutex
	buf  bytes.Buffer
}

func (t *lineTagger) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buf.Write(p)
	for {
		line, err := t.buf.ReadBytes('\n')
		if err != nil {
			// No newline yet — push back and stop.
			t.buf.Reset()
			t.buf.Write(line)
			break
		}
		if _, err := t.out.Write([]byte("[" + t.name + "] ")); err != nil {
			return 0, err
		}
		if _, err := t.out.Write(line); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
