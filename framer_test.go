package dcnet

import (
	"io"
	"io/ioutil"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// Interface assertions
var _ Framer = (*RTPFramer)(nil)
var _ FrameReader = (*RTPFrameReader)(nil)
var _ FrameWriter = (*RTPFrameWriter)(nil)

func TestRTPFramer(t *testing.T) {

	testCases := []struct {
		sequence string
	}{
		{""},
		{"abc"},
		{RandString(65535)},
	}

	for i, testCase := range testCases {
		// TODO: make the pipe split the message
		pr, pw := io.Pipe()

		r, err := NewRTPFrameReader(pr)
		if err != nil {
			t.Fatalf("failed to create frame reader: %v", err)
		}
		defer r.Close()

		t.Run("NewRTPFramer_"+strconv.Itoa(i), func(t *testing.T) {

			w, err := NewRTPFrameWriter(len(testCase.sequence), pw)
			if err != nil {
				t.Fatalf("failed to create frame writer: %v", err)
			}
			defer w.Close()

			done := make(chan struct{})

			time.AfterFunc(time.Second*2, func() {
				// Avoid extreme waiting time on blocking bugs
				panic("timeout")
			})

			go func() {
				defer close(done)
				i, err := w.Write([]byte(testCase.sequence))
				if i != len(testCase.sequence) {
					t.Fatalf("short write")
				}
				if err != nil {
					t.Fatalf("failed to write: %v", err)
				}
			}()

			result, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatalf("failed to read: %v", err)
			}

			<-done

			if string(result) != testCase.sequence {
				t.Errorf(string(result) + " != " + testCase.sequence)
			}

		})
	}

}

func RandString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
