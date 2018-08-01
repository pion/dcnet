package dcnet

// NilAddr is an empty address.
// It can be used in scenario's where multiple
// peers are not needed.
type NilAddr struct {
	ID string
}

func (a *NilAddr) Network() string {
	return "WebRTC"
}

func (a *NilAddr) String() string {
	return ""
}

// IDAddr identifies a peer by a simple ID.
type IDAddr struct {
	ID string
}

func (a *IDAddr) Network() string {
	return "WebRTC"
}

func (a *IDAddr) String() string {
	return a.ID
}

// SessionAddr identifies a peer by a
// combination of APIKey, RoomID and Session ID.
type SessionAddr struct {
	APIKey    string
	RoomID    string
	SessionID string
}

func (a *SessionAddr) Network() string {
	return "WebRTC"
}

func (a *SessionAddr) String() string {
	return a.APIKey +
		"/" + a.RoomID +
		"/" + a.SessionID
}
