package dcnet

import (
	"encoding/binary"
	"errors"
	"io"
)

// Framer allows sending and receiving packets with arbitrary over a io.ReadWriteCloser
type Framer interface {
	io.ReadWriteCloser
}

// Framer allows receiving packets with arbitrary length over a io.ReadCloser
type FrameReader interface {
	io.ReadCloser
}

// Framer allows sending packets with arbitrary length over a io.WriteCloser
type FrameWriter interface {
	io.WriteCloser
}

// RTPFramer implements a RFC4571 Framer
type RTPFramer struct {
	r *RTPFrameReader
	w *RTPFrameWriter
}

func NewRTPFramer(length int, c io.ReadWriteCloser) (*RTPFramer, error) {
	r, err := NewRTPFrameReader(c)
	if err != nil {
		return nil, err
	}
	w, err := NewRTPFrameWriter(length, c)
	if err != nil {
		return nil, err
	}
	return &RTPFramer{
		r: r,
		w: w,
	}, nil
}

func (f *RTPFramer) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}

func (f *RTPFramer) Write(p []byte) (n int, err error) {
	return f.w.Write(p)
}

func (f *RTPFramer) Close() error {
	return f.r.Close()
}

// RTPFrameReader implements a RFC4571 FrameReader
// This framer only works over reliable, ordered connections
type RTPFrameReader struct {
	length    uint16
	inFrame   bool
	curLength uint16
	c         io.ReadCloser
}

// NewRTPFrameReader allows sending packets of up to 65535 bytes over the ReadWriteCloser
func NewRTPFrameReader(c io.ReadCloser) (*RTPFrameReader, error) {
	return &RTPFrameReader{
		c: c,
	}, nil
}

// Read reads a packet from the underlying stream. It returns EOF
func (f *RTPFrameReader) Read(p []byte) (n int, err error) {
	n, err = f.c.Read(p)
	if err != nil {
		return n, err
	}

	if !f.inFrame {
		f.length = binary.BigEndian.Uint16(p[:2])
		p = p[2:]
		n -= 2
		f.inFrame = true
	}

	newLen := int(f.curLength) + n
	if newLen > int(f.length) {
		return 0, errors.New("Receiving packet too long")
	}

	f.curLength = uint16(newLen)

	// Close out the packet
	if f.length == f.curLength {
		f.inFrame = false
		f.curLength = 0
		return n, io.EOF // TODO: Is this fine or should we return io.EOF forever and re-create the RTPFrameReader?
	}

	return n, nil
}

// Close closes the underlying ReadCloser
func (f *RTPFrameReader) Close() error {
	return f.c.Close()
}

// RTPFrameWriter implements a RFC4571 FrameWriter
// This framer only works over reliable, ordered connections
type RTPFrameWriter struct {
	length    uint16
	inFrame   bool
	curLength uint16
	c         io.WriteCloser
}

// NewRTPFrameWriter allows sending packets of up to 65535 bytes over the ReadWriteCloser
func NewRTPFrameWriter(length int, c io.WriteCloser) (*RTPFrameWriter, error) {
	if length > 65535 {
		return nil, errors.New("Maximum length (65535) exeeded")
	}
	return &RTPFrameWriter{
		length: uint16(length),
		c:      c,
	}, nil
}

func (f *RTPFrameWriter) Write(p []byte) (n int, err error) {
	if int(f.curLength)+len(p) > int(f.length) {
		return 0, errors.New("Sending packet too long")
	}

	if !f.inFrame {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, f.length)
		_, err := f.c.Write(b) // Assumes we can at least write 2 bytes (or get an error)
		if err != nil {
			return 0, err
		}
		f.inFrame = true
	}

	if len(p) > 0 {
		n, err = f.c.Write(p)
	}

	f.curLength = f.curLength + uint16(n)

	// Close out the packet
	if f.length == f.curLength {
		f.curLength = 0
		f.inFrame = false
	}

	return n, err
}

// Close closes the underlying WriteCloser
func (f *RTPFrameWriter) Close() error {
	return f.c.Close()
}
