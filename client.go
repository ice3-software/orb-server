
//
// A primative type that models an Orb.
//
type Orb struct {
	X			float32
	Y 			float32
	ID			string
}


//
// A construct responsibile for reading and writing Orb changes between the
// TCP client and server.
//
type OrbClient struct {

	//
	// The client's Orb data.
	//
	Orb				Orb

	//
	// The client'd TCP connection.
	//
	Conn 			*net.TCPConn

	//
	// Client to broadcast Orb changes to its Room
	//
	Read			chan Orb

	//
	// For room to broadcast Orb changes to its clients
	//
	Write			chan Orb

	//
	// For client to notify the room that it has been disconnected
	//
	Disconnect 	chan bool

}

func (self *OrbClient) broadcastDisconnect() {
	self.Disconnect <-true
}

func (self *OrbClient) write() {

	for {

		changedOrb := <-self.Write

		// TODO: Serialise orb
		// _, err := self.Conn.Write(...)
		// if err != nil {
		//		self.broadcastDisconnect()
		//		return
		// }
	}

}

func (self *OrbClient) read() {

	for {

		msgBuf := make([]byte, 2048)
		msgLen, err := self.Conn.Read(msgBuf)

		if err == io.EOF {
			self.broadcastDisconnect()
			return
		} else {
			// TODO: Parse message into the orb model properly
			self.Orb = Orb{
				X: 	123
				Y: 	123
				ID: self.Orb.ID
			}
			self.Read <-self.Orb
		}
	}
}

func (self *OrbClient) Close() {
	self.Conn.Close()
}

func NewOrbClient(conn *net.TCPConn) *OrbClient {

	orb := &OrbClient{
		Orb: Orb{
			ID: "123", // TODO: Make this a unique string
		},
		Conn: conn,
		Write: make(chan Orb),
		Read: make(chan Orb),
		Disconnect: make(chan bool),
	}

	go orb.read()
	go orb.write()

	return orb
}
