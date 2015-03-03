
package main

import (
	"net"
	"fmt"
	//"sync"
)

type JoinRequest struct {
	Conn	net.Conn
	Joined	chan bool
}

func (self JoinRequest) broadcastJoined() {

	// TODO: Is it safe to check whether a channel is nil? Are channels
	// mutatable across different threads?

	if self.Joined != nil {
		self.Joined <-true
		close(self.Joined)
	}
}

type Room struct {
	clients		[]*OrbClient
	join		chan JoinRequest
	sharedRead 	chan Orb
	publicRead 	chan Orb
}

func (self *Room) Join() chan<- JoinRequest {
	return self.join
}

func (self *Room) Read() <-chan Orb {
	return self.publicRead
}

func (self *Room) broadcastOrb(orb Orb) {
	for _, client := range self.clients {
		fmt.Println("Printing %q to %q", orb, client)
		client.Write() <-orb
	}
}

func (self *Room) mainLoop() {

	self.clients = make([]*OrbClient, 0, 5)
	self.sharedRead = make(chan Orb)

	for {
		select {

			case joinReq := <-self.join:
			newClient := NewOrbClient(joinReq.Conn, self.sharedRead)
			self.clients = append(self.clients, newClient)
			joinReq.broadcastJoined()
			break

			case orbChange := <-self.sharedRead:
			self.broadcastOrb(orbChange)
			break
		}
	}
}

func NewRoom() *Room {
	room := &Room{
		join: make(chan JoinRequest),
	}
	go room.mainLoop()
	return room
}
