package dcnet

import "net"

// Listener creates an io.listener around an existing RTCPeerConnection
// It is up to the supplier of the RTCPeerConnection to
type Listener struct {
	s Signaler
}

func NewListener(s Signaler) (*Listener, error) {
	res := &Listener{
		s: s,
	}

	return res, nil
}

func (l *Listener) Accept() (net.Conn, error) {
	dc, raddr, err := l.s.Accept()
	if err != nil {
		return nil, err
	}

	conn, err := NewConn(dc, l.s.Addr(), raddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (l *Listener) Close() error {
	return l.s.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.s.Addr()
}
