<h1 align="center">
  <br>
  Pion DCNet
  <br>
</h1>
<h4 align="center">Net interfaces around WebRTC data channels.</h4>
<br>

**Warning**: This package is a proof of concept. Feel free to play around, log issues or even contribute but don't rely on anything!

Package DCNet augments the net.* interfaces over a WebRTC data channels.

### Usage
The idea of the package is that you implement a ``Signaler`` that negotiates data channels and get a net.Listener and net.Conn on top of them. The following pseudo-code shows how this works: 

Listening side:
``` Go
signaler := NewSignaler(someOptions)
listener := dcnet.NewListener(signaler)
for {
	conn, err := listener.Accept()
	check(err)
	_, err = conn.Write([]byte("Hallo"))
	check(err)
}
```

Dialing side:
``` Go
signaler := NewSignaler(someOptions)
conn, err := dcnet.Dial(signaler)
check(err)

msg := make([]byte, 100)
_, err = conn.Read(msg)
check(err)
fmt.Println(string(msg))
```
### Examples
Please refer to the examples directory for working examples.