package conn

import (
	"bytes"
	"net"
)

// Recorder is a connection wrapper that can record all bytes read and play them
// back again.  Reset() can be called at any point, after which Read() will
// start returning bytes from the start of the stream again.  StopRecording()
// permanently prevents further recording.
type Recorder struct {
	*net.TCPConn

	buf      bytes.Buffer
	reread   bool
	norecord bool
	i        int
}

func (c *Recorder) Read(b []byte) (int, error) {
	if c.reread {
		if c.i < c.buf.Len() {
			n := copy(b, c.buf.Bytes()[c.i:])
			c.i += n
			return n, nil
		}
		c.reread = false
	}

	if c.norecord {
		c.buf = bytes.Buffer{}
	}

	n, err := c.TCPConn.Read(b)
	if !c.norecord {
		c.buf.Write(b[:n])
	}
	return n, err
}

// Rewind causes Read() to start returning bytes from the start of the stream
// again.
func (c *Recorder) Rewind() {
	if c.norecord {
		panic("Rewind() may not be called after StopRecording() has been called")
	}

	c.i = 0
	c.reread = true
}

// StopRecording permanently prevents further recording.
func (c *Recorder) StopRecording() {
	c.norecord = true
}
