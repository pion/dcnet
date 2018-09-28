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

// Signaler
type Signaler interface {
	Dial() (*webrtc.RTCDataChannel, net.Addr, error)
	Accept() (*webrtc.RTCDataChannel, net.Addr, error)
	Close() error
	Addr() net.Addr
}

// RWSignaler is a simple signaler over an io.ReadWriteCloser
// It plainly exchanges the SDP offer/answer without peer address.
type RWSignaler struct {
	c      io.ReadWriteCloser
	config webrtc.RTCConfiguration // Maybe make this global?
}

// NewRWSignaler creates a new RWSignaler
func NewRWSignaler(c io.ReadWriteCloser, rtcconfig webrtc.RTCConfiguration) *RWSignaler {
	s := &RWSignaler{
		c:      c,
		config: rtcconfig,
	}

	return s
}

// Accept creates WebRTC DataChannels by signaling over the ReadWriteCloser
func (r *RWSignaler) Dial() (*webrtc.RTCDataChannel, net.Addr, error) {
	return r.doSignal(true)
}

// Accept creates WebRTC DataChannels by signaling over the ReadWriteCloser
func (r *RWSignaler) Accept() (*webrtc.RTCDataChannel, net.Addr, error) {
	return r.doSignal(false)
}

func (r *RWSignaler) doSignal(init bool) (*webrtc.RTCDataChannel, net.Addr, error) {
	c, err := webrtc.New(r.config)
	if err != nil {
		return nil, nil, err
	}

	var dc *webrtc.RTCDataChannel
	var addr net.Addr

	if init {
		dc, err = c.CreateDataChannel("data", nil)
		if err != nil {
			return nil, nil, err
		}

		// TODO: migrate to OnNegotiationNeeded when available
		offer, err := c.CreateOffer(nil)
		if err != nil {
			return nil, nil, err
		}

		b, err := json.Marshal(offer)
		if err != nil {
			return nil, nil, err
		}

		fmt.Println("Sending initial offer", string(b))
		fmt.Println()
		f, err := NewRTPFrameWriter(len(b), r.c)
		if err != nil {
			return nil, nil, err
		}

		// TODO: Don't assume we can write the entire offer at once
		_, err = f.Write(b)
		if err != nil {
			return nil, nil, err
		}
		addr = &NilAddr{}
	}

	go func() {
		for {
			f, err := NewRTPFrameReader(r.c)
			if err != nil {
				// TODO: Return error from Accept()
				log.Println(err)
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				// TODO: Return error from Accept()
				log.Println(err)
			}

			fmt.Println()
			var desc webrtc.RTCSessionDescription
			err = json.Unmarshal(b, &desc)
			if err != nil {
				// TODO: Return error from Accept()
				log.Println(err)
			}

			if err := c.SetRemoteDescription(desc); err != nil {
				panic(err)
			}
			fmt.Println("Got", desc.Type, string(b))
			log.Println(desc.Type, webrtc.RTCSdpTypeOffer, desc.Type == webrtc.RTCSdpTypeOffer)
			if desc.Type == webrtc.RTCSdpTypeOffer {
				// Sets the LocalDescription, and starts our UDP listeners
				answer, err := c.CreateAnswer(nil)
				if err != nil {
					panic(err)
				}

				b, err := json.Marshal(answer)
				if err != nil {
					// TODO: Return error from Accept()
					log.Println(err)
				}

				fmt.Println("Sending answer", string(b))
				fmt.Println()
				f, err := NewRTPFrameWriter(len(b), r.c)
				if err != nil {
					// TODO: Return error from Accept()
					log.Println(err)
				}

				// TODO: Don't assume we can write the entire answer at once
				_, err = f.Write(b)
				if err != nil {
					// TODO: Return error from Accept()
					log.Println(err)
				}
			}
		}
	}()

	if dc == nil {
		res := make(chan struct {
			d *webrtc.RTCDataChannel
			a net.Addr
		})

		c.OnDataChannel = func(d *webrtc.RTCDataChannel) {
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
		dc = e.d
		addr = e.a
	}

	return dc, addr, nil
}

// Close closed the ReadWriteCloser
func (r *RWSignaler) Close() error {
	return r.c.Close()
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
	panic("TODO")
	return nil, nil, nil
}

func (s *MultiSignaler) Close() error {
	var closeErr error
	for _, signaler := range s.s {
		err := signaler.Close()
		if err != nil {
			closeErr = err
		}
	}
	return closeErr
}

func (s *MultiSignaler) Addr() net.Addr {
	// We assume all signalers use the same addr
	return s.s[0].Addr()
}
