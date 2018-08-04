package dcnet

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/pions/webrtc"
)

// Addr + Signaler -> DataChannel
// Signaler
type Signaler interface {
	Accept() (*webrtc.RTCDataChannel, net.Addr, error)
	Close() error
	Addr() net.Addr
}

// RWSignaler is a simple signaler over an io.ReadWriteCloser
// It plainy exchanges the SDP offer/answer without peer address.
type RWSignaler struct {
	c        io.ReadWriteCloser
	config   webrtc.RTCPeerConfiguration // Maybe make this global?
	initiate bool
}

// NewRWSignaler creates a new RWSignaler
func NewRWSignaler(c io.ReadWriteCloser, rtcconfig webrtc.RTCPeerConfiguration, initiate bool) (*RWSignaler, error) {
	s := &RWSignaler{
		c:      c,
		config: rtcconfig,
		initiate, initiate,
	}

	return s, nil
}

// Accept creates WebRTC DataChannels by signaling over the ReadWriteCloser
// TODO: Abstract further?
func (r *RWSignaler) Accept() (*webrtc.RTCDataChannel, net.Addr, error) {
	c, err := webrtc.New(r.config)
	if err != nil {
		return nil, nil, err
	}

	var dc *webrtc.DataC

	if r.initiate {
		dc, err := c.CreateDataChannel("data", nil)
		if err != nil {
			return nil, nil, err
		}

		offer, err := c.CreateOffer(nil)
		if err != nil {
			return nil, nil, err
		}

		b, err := json.Marshal(offer)
		if err != nil {
			return err
		}

		// TODO: Don't assume we can write the entire offer at once
		_, err = r.c.Write(b)
		if err != nil {
			return nil, nil, err
		}
	}

	go func() {
		for {
			b, err := ioutil.ReadAll(r.c) // Blocks forever?
			if err != nil {
				// TODO: Return error from Accept()
				log.Println(err)
			}

			var desc webrtc.RTCSessionDescription
			err = json.Unmarshal(b, &desc)
			if err != nil {
				// TODO: Return error from Accept()
				log.Println(err)
			}

			if err := peerConnection.SetRemoteDescription(desc); err != nil {
				panic(err)
			}
			if desc.Type == webrtc.RTCSdpTypeOffer {
				// Sets the LocalDescription, and starts our UDP listeners
				answer, err := peerConnection.CreateAnswer(nil)
				if err != nil {
					panic(err)
				}
				_, err = r.c.Write(p)
				if err != nil {
					// TODO: Return error from Accept()
					log.Println(err)
				}
			}
		}
	}()

	res := make(chan struct {
		d *webrtc.RTCDataChannel
		a net.Addr
	})

	c.Ondatachannel = func(d *webrtc.RTCDataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label, d.ID)
		res <- struct {
			d *webrtc.RTCDataChannel
			a net.Addr
		}{
			d: d,
			a: &NilAddr{},
		}
	}

	e := <-res

	return e.d, e.a, nil
}

// Close closed the ReadWriteCloser
func (r *RWSignaler) Close() error {
	return c.Close()
}

func (r *RWSignaler) Addr() net.Addr {
	return &NilAddr{}
}

// MultiSignaler combines the power of many signalers
type MultiSignaler struct {
	s []Signaler
}

// NewMultiSignaler creates a new MultiSignaler
func NewMultiSignaler(set ...Signaler) (*MultiSignaler, error) {
	res := &MultiSignaler{
		s: set,
	}
	return res, nil
}

func (s *MultiSignaler) Accept() (*webrtc.RTCDataChannel, net.Addr, error) {
	// TODO: accept on all signalers
}

func (s *MultiSignaler) Close() error {
	var closeerr error
	for _, signaler := range s.s {
		err := signaler.Close()
		if err != nil {
			closeerr = err
		}
	}
	return closeerr
}

func (s *MultiSignaler) Addr() net.Addr {
	// We assume all signalers use the same addr
	return s[0].Addr()
}
