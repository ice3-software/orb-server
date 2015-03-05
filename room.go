
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
	started 	chan bool
}

func (self *Room) Started() chan<- started {
	return self.started
}

func (self *Room) Join() chan<- JoinRequest {
	return self.join
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
	self.started = make(chan bool)

	var started bool

	for {

		var read chan Orb

		if started {
			read = self.sharedRead
		}

		select {

			case joinReq := <-self.join:
			fmt.Printf("New orb client connected. %q clients\n", len(self.clients))
			newClient := NewOrbClient(joinReq.Conn, self.sharedRead)
			self.clients = append(self.clients, newClient)
			joinReq.broadcastJoined()

			if len(self.clients) > 5 {
				started = true
				self.started <-started
			}

			break

			case orbChange := <-read:
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
