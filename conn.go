package dcnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/datachannel"
)

// NewConn creates a Conn around a data channel. The data channel is assumed to be open already.
func NewConn(dc *webrtc.RTCDataChannel, laddr net.Addr, raddr net.Addr) (net.Conn, error) {
	r, w := io.Pipe()

	res := &Conn{
		dc:    dc,
		laddr: laddr,
		raddr: raddr,
		p:     r,
	}

	go func() {
		dc.Lock()
		defer dc.Unlock()
		dc.Onmessage = func(payload datachannel.Payload) {
			switch p := payload.(type) {
			// case *datachannel.PayloadString:
			// 	fmt.Printf("Message '%s' from DataChannel '%s' payload '%s'\n", p.PayloadType().String(), d.Label, string(p.Data))
			case *datachannel.PayloadBinary:
				w.Write(p.Data)
				// default:
				// 	fmt.Printf("Message '%s' from DataChannel '%s' no payload \n", p.PayloadType().String(), d.Label)
			}
		}
	}()

	return res, nil
}

// Conn is a net.Conn over a datachannel
type Conn struct {
	dc       *webrtc.RTCDataChannel
	laddr    net.Addr
	raddr    net.Addr
	p        *io.PipeReader
	isClosed bool
}

// Read reads data from the underlying the data channel
func (c *Conn) Read(b []byte) (int, error) {
	if c.isClosed {
		return 0, errors.New("read on closed conn")
	}
	i, err := c.p.Read(b)
	return i, err
}

// Write writes the data to the underlying data channel
func (c *Conn) Write(b []byte) (int, error) {
	if c.isClosed {
		return 0, errors.New("write on closed conn")
	}
	err := c.dc.Send(datachannel.PayloadBinary{Data: b})
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// Close closes the datachannel and peerconnection
func (c *Conn) Close() error {

	if c.isClosed {
		return errors.New("close on closed conn")
	}

	// TODO: Locking
	c.isClosed = true

	// Unblock readers
	err := c.p.Close()
	if err != nil {
		fmt.Println("failed to close pipe:", err)
		return err
	}

	// TODO: Implement datachannel closing procedure
	// c.dc.Close()
	// TODO: cleanup peerconnection
	return nil
}

func (c *Conn) LocalAddr() net.Addr {
	return c.laddr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.raddr
}

// SetDeadline
func (c *Conn) SetDeadline(t time.Time) error {
	panic("TODO")
	return nil
}

// SetReadDeadline
func (c *Conn) SetReadDeadline(t time.Time) error {
	panic("TODO")
	return nil

}

// SetWriteDeadline
func (c *Conn) SetWriteDeadline(t time.Time) error {
	panic("TODO")
	return nil
}
