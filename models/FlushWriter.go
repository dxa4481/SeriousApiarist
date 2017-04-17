package models

import (
	"io"
	"net/http"
)

// FlushWriter can be written to and flushed to the HTTP stream
type FlushWriter struct {
	f http.Flusher
	W io.Writer
}

// NewFlushWriter allows you to write to streaming responses
func NewFlushWriter(w http.ResponseWriter) FlushWriter {
	fw := FlushWriter{W: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}
	return fw
}

func (fw *FlushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.W.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}
