package dcnet

import (
	"net"
	"testing"
	"time"

	"github.com/pions/webrtc"
)

// Interface assertions
var _ Signaler = (*RWSignaler)(nil)
var _ Signaler = (*MultiSignaler)(nil)

func TestRWSignaler(t *testing.T) {

	// Avoid extreme waiting time on blocking bugs
	time.AfterFunc(time.Second*2, func() {
		panic("timeout")
	})

	connA, connB := net.Pipe()

	sigA := NewRWSignaler(connA, webrtc.RTCConfiguration{}, true)
	sigB := NewRWSignaler(connB, webrtc.RTCConfiguration{}, false)

	doneA := validateAccept(t, sigA)
	doneB := validateAccept(t, sigB)

	<-doneA
	<-doneB
}

func validateAccept(t *testing.T, sig Signaler) chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		_, _, err := sig.Accept()
		if err != nil {
			t.Errorf("failed to accept: %v", err)
			return
		}

		// TODO: test sending something once pions to pions works
	}()

	return done
}
