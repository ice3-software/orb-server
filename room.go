
package main

import (
	"net"
	//"sync"
)

type Room struct {
	clients		[]*OrbClient
	join		chan net.Conn
	sharedRead 	chan Orb
}

func (self *Room) Join() chan<- net.Conn{
	return self.join
}

func (self *Room) broadcastOrb(orb Orb) {
	for _, client := range self.clients {
		client.Write() <-orb
	}
}

func (self *Room) mainLoop() {

	self.clients = make([]*OrbClient, 0, 5)
	self.sharedRead = make(chan Orb)

	for {
		select {

			case conn := <-self.join:
			newClient := NewOrbClient(conn, self.sharedRead)
			self.clients = append(self.clients, newClient)
			break

			case orbChange := <-self.sharedRead:
			self.broadcastOrb(orbChange)
			break
		}
	}
}

func NewRoom() *Room {
	room := &Room{}
	go room.mainLoop()
	return room
}
