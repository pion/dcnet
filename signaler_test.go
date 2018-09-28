package dcnet

import (
	"log"
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
	lim := time.AfterFunc(time.Second*5, func() {
		panic("timeout")
	})
	defer lim.Stop()

	connA, connB := net.Pipe()

	sigA := NewRWSignaler(connA, webrtc.RTCConfiguration{}, true)
	sigB := NewRWSignaler(connB, webrtc.RTCConfiguration{}, false)

	doneA := validateAccept(t, "A", sigA)
	doneB := validateAccept(t, "B", sigB)

	<-doneA
	<-doneB
}

func validateAccept(t *testing.T, name string, sig Signaler) chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		_, _, err := sig.Accept()
		if err != nil {
			t.Errorf("failed to accept on %s: %v", name, err)
			return
		}
		log.Println("Accept on", name)

	}()

	return done
}
