package dcnet

import "net"

// Interface assertions
var _ net.Addr = (*NilAddr)(nil)
var _ net.Addr = (*IDAddr)(nil)
var _ net.Addr = (*SessionAddr)(nil)
