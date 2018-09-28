package dcnet

import (
	"net"
)

func Dial(signaler Signaler) (net.Conn, error) {
	dc, raddr, err := signaler.Dial()
	if err != nil {
		return nil, err
	}

	// Ensure the channel is open
	ensureOpen(dc)

	conn, err := NewConn(dc, signaler.Addr(), raddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
