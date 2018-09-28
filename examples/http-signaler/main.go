package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/pions/dcnet"
	"github.com/pions/webrtc"
)

func main() {
	init := flag.Bool("init", false, "initiate the connection")
	addr := flag.String("address", ":50000", "signaler address")
	flag.Parse()

	config := webrtc.RTCConfiguration{}
	signaler := NewHTTPSignaler(config, *addr)

	var conn net.Conn
	var err error
	if *init {
		conn, err = dcnet.Dial(signaler)
		check(err)
	} else {
		listener := dcnet.NewListener(signaler)
		conn, err = listener.Accept()
		check(err)
	}

	// Send a message to the other side
	_, err = conn.Write([]byte("Hello world!"))
	check(err)
	fmt.Println("Sent message")

	// Receive the message from the other side
	b := make([]byte, 100)
	_, err = conn.Read(b)
	check(err)
	fmt.Println("Received message:", string(b))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
