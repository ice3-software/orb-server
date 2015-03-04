
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
	fmt.Println("Broadcasting change")
	for index, client := range self.clients {
		fmt.Printf("Printing {%q, %f, %f} to client %d\n", orb.ID, orb.X, orb.Y, index)
		client.Write() <-orb
	}
}

func (self *Room) mainLoop() {

	self.clients = make([]*OrbClient, 0, 5)
	self.sharedRead = make(chan Orb)

	for {
		select {

			case joinReq := <-self.join:
			fmt.Printf("New orb client connected. %q clients\n", len(self.clients))
			newClient := NewOrbClient(joinReq.Conn, self.sharedRead)
			self.clients = append(self.clients, newClient)
			joinReq.broadcastJoined()
			break

			case orbChange := <-self.sharedRead:
			fmt.Println("Orb changed: ", orbChange.ID)
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
