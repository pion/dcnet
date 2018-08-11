package dcnet

import "net"

// Interface assertions
var _ net.Conn = (*Conn)(nil)
