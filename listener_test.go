package dcnet

import "net"

// Interface assertions
var _ net.Listener = (*Listener)(nil)
