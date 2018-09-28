package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/pions/dcnet"
	"github.com/pions/webrtc"
)

func NewHTTPSignaler(config webrtc.RTCConfiguration, address string) *HTTPSignaler {
	return &HTTPSignaler{
		config:  webrtc.RTCConfiguration{},
		address: address,
	}
}

type HTTPSignaler struct {
	config  webrtc.RTCConfiguration
	address string
}

func (r *HTTPSignaler) Dial() (*webrtc.RTCDataChannel, net.Addr, error) {
	c, err := webrtc.New(r.config)
	if err != nil {
		return nil, nil, err
	}

	var dc *webrtc.RTCDataChannel

	dc, err = c.CreateDataChannel("data", nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: migrate to OnNegotiationNeeded when available
	offer, err := c.CreateOffer(nil)
	if err != nil {
		return nil, nil, err
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(offer)
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.Post("http://"+r.address, "application/json; charset=utf-8", b)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var answer webrtc.RTCSessionDescription
	err = json.NewDecoder(resp.Body).Decode(&answer)
	if err != nil {
		return nil, nil, err
	}

	if err := c.SetRemoteDescription(answer); err != nil {
		return nil, nil, err
	}
	return dc, &dcnet.NilAddr{}, nil
}

func (r *HTTPSignaler) Accept() (*webrtc.RTCDataChannel, net.Addr, error) {
	c, err := webrtc.New(r.config)
	if err != nil {
		return nil, nil, err
	}

	var dc *webrtc.RTCDataChannel

	res := make(chan *webrtc.RTCDataChannel)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var offer webrtc.RTCSessionDescription
		err := json.NewDecoder(r.Body).Decode(&offer)
		if err != nil {
			log.Println("Failed to decode offer:", err)
			return
		}

		if err := c.SetRemoteDescription(offer); err != nil {
			log.Println("Failed to set remote description:", err)
			return
		}

		answer, err := c.CreateAnswer(nil)
		if err != nil {
			log.Println("Failed to create answer:", err)
			return
		}

		err = json.NewEncoder(w).Encode(answer)
		if err != nil {
			log.Println("Failed to encode answer:", err)
			return
		}

		c.OnDataChannel = func(d *webrtc.RTCDataChannel) {
			res <- d
		}

	})

	go http.ListenAndServe(r.address, nil)

	dc = <-res
	return dc, &dcnet.NilAddr{}, nil
}

func (r *HTTPSignaler) Close() error {
	return nil
}

func (r *HTTPSignaler) Addr() net.Addr {
	return &dcnet.NilAddr{}
}
