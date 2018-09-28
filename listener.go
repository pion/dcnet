package dcnet

import (
	"net"

	"github.com/pions/webrtc"
)

// Listener creates an io.listener around an existing RTCPeerConnection
// It is up to the supplier of the RTCPeerConnection to
type Listener struct {
	s    Signaler
	addr net.Addr
}

func NewListener(s Signaler) *Listener {
	res := &Listener{
		s: s,
	}

	return res
}

func (l *Listener) Accept() (net.Conn, error) {
	dc, raddr, err := l.s.Accept()
	if err != nil {
		return nil, err
	}

	// Ensure the channel is open
	ensureOpen(dc)

	conn, err := NewConn(dc, l.s.Addr(), raddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ensureOpen(dc *webrtc.RTCDataChannel) {
	done := make(chan struct{})
	open := func() {
		select {
		case <-done:
		default:
			close(done)
		}
	}
	dc.OnOpen = open

	if dc.ReadyState == webrtc.RTCDataChannelStateOpen {
		open()
	}
	<-done
}

func (l *Listener) Close() error {
	return l.s.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.s.Addr()
}
