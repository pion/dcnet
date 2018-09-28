package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pions/dcnet"
	"github.com/pions/webrtc"
)

func main() {

	// Avoid extreme waiting time on blocking bugs
	lim := time.AfterFunc(time.Second*10, func() {
		panic("timeout")
	})
	defer lim.Stop()

	connA, connB := net.Pipe()

	sigA := dcnet.NewRWSignaler(connA, webrtc.RTCConfiguration{}, true)
	sigB := dcnet.NewRWSignaler(connB, webrtc.RTCConfiguration{}, false)

	doneA := validateAccept("A", sigA)
	doneB := validateAccept("B", sigB)

	<-doneA
	<-doneB

	log.Println("All good!")

}

func validateAccept(name string, sig dcnet.Signaler) chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		_, _, err := sig.Accept()
		if err != nil {
			panic(fmt.Sprintf("failed to accept on %s: %v", name, err))
			return
		}
		log.Println("Accept on", name)

		// TODO: test sending something once pions to pions works
	}()

	return done
}
