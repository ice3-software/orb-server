
package main

import (
	"code.google.com/p/go-uuid/uuid"
	"gopkg.in/mgo.v2/bson"
	"net"
	"io"
	"fmt"
)

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
	orb				Orb

	//
	// The client'd connection. This will be an *net.TCPConn. We're using
	// an interface here for mocking purposes.
	//
	conn 			net.Conn

	//
	// Client to broadcast Orb changes to its Room
	//
	read			chan Orb

	//
	// For room to broadcast Orb changes to its clients
	//
	write			chan Orb

	//
	// For client to notify the room that it has been disconnected
	//
	disconnect 		chan bool

}

//
// Getters
//

func (self *OrbClient) Orb() Orb {
	return self.orb
}

func (self *OrbClient) Conn() net.Conn {
	return self.conn
}

func (self *OrbClient) Read() <-chan Orb {
	return self.read
}

func (self *OrbClient) Write() chan<- Orb {
	return self.write
}

func (self *OrbClient) Disconnect() <-chan bool {
	return self.disconnect
}

//
// Unexported methods
//

func (self *OrbClient) broadcastDisconnect() {
	self.disconnect <-true
}

func (self *OrbClient) writeLoop() {

	for {

		changedOrb := <-self.write
		bsonBuf, err := bson.Marshal(changedOrb)

		if err != nil {
			// TODO Handle error...
			fmt.Println("Error marshalling: ", err)
		}
		_, err = self.conn.Write(bsonBuf)

		if err == io.EOF {
			self.broadcastDisconnect()
			return
		} else if err != nil {
			// TODO Handle error...
			fmt.Println("Error writing: ", err)
		}

	}

}

func (self *OrbClient) readLoop() {

	for {

		bsonBuf := make([]byte, 2048)
		_, err := self.conn.Read(bsonBuf)

		if err == io.EOF {
			self.broadcastDisconnect()
			return
		} else if err == nil {

			newOrb := &Orb{}
			if err := bson.Unmarshal(bsonBuf, newOrb); err != nil {
				// TODO Handle error...
				fmt.Println("Error unmarshalling: ", err)
			}
			newOrb.ID = self.orb.ID

			self.orb = *newOrb
			self.read <-self.orb

		} else {
			// TODO Handle error...
			fmt.Println("Error reading: ", err)
		}
	}
}

func (self *OrbClient) Close() error {
	return self.conn.Close()
}

//
// Ctor
//

func NewOrbClient(conn net.Conn) *OrbClient {

	client := &OrbClient{
		orb: Orb{
			ID: uuid.New(),
		},
		conn: conn,
		write: make(chan Orb),
		read: make(chan Orb),
		disconnect: make(chan bool),
	}

	go client.readLoop()
	go client.writeLoop()

	return client
}
