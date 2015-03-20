
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
	limit		int
}

func (self *Room) Started() <-chan bool {
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

	join := self.join
	var startedCh chan bool
	var started bool

	for {

		// Allow reads if the room has actually started, otherwise turn off
		// that select case with a nil channel

		var read chan Orb

		if started {
			read = self.sharedRead
		}

		select {

			case startedCh <-started :
				fmt.Println("Sent starting message...")
				//for _, client := range self.clients {
				//	client.Listen()
				//}
				startedCh = nil
				break

			case joinReq := <-join :

				newClient := NewOrbClient(joinReq.Conn, self.sharedRead)
				self.clients = append(self.clients, newClient)
				joinReq.broadcastJoined()

				fmt.Println("Joined: ", len(self.clients))

				if len(self.clients) >= self.limit {

					fmt.Println("Hit limit, starting...")

					// Open the started started channel so that we can send 1 message down. The idea
					// here is that the join channel will be closed but we will allow client changes
					// to be broadcasted before someone has recieved a message on the Started chan.

					started = true
					startedCh = self.started

					// Nil the internal join channel so we don't accept any more join requests.
					// Any routines that have a handle on this channel will just block when
					// sending requests. They should listen for messages on the Started chan and
					// subsequently assume that this room is now closed.

					join = nil
				}

				break

			case orbChange := <-read :

				fmt.Println("Orb changed: ", orbChange.ID)
				self.broadcastOrb(orbChange)
				break

		}
	}
}

func NewRoom(limit int) *Room {
	room := &Room{
		limit: limit,
		join: make(chan JoinRequest),
		started: make(chan bool),
	}
	go room.mainLoop()
	return room
}
