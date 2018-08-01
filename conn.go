package dcnet

import (
	"io"
	"net"
	"time"

	"github.com/pions/webrtc"
)

func NewConn(dc *webrtc.RTCDataChannel, laddr net.Addr, raddr net.Addr) (net.Conn, error) {
	res := &Conn{
		dc:    dc,
		laddr: laddr,
		raddr: raddr,
	}

	// TODO: setup some copying

	return res, nil
}

type Conn struct {
	dc    *webrtc.RTCDataChannel
	laddr net.Addr
	raddr net.Addr
	p     *io.PipeReader
}

func (c *Conn) Read(b []byte) (int, error) {

}

func (c *Conn) Write(b []byte) (n int, err error) {

}

func (c *Conn) Close() error {

}

func (c *Conn) LocalAddr() net.Addr {
	return c.laddr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.raddr
}

func (c *Conn) SetDeadline(t time.Time) error {

}

func (c *Conn) SetReadDeadline(t time.Time) error {

}

func (c *Conn) SetWriteDeadline(t time.Time) error {

}
