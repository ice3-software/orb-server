
package main

import (
	"net"
)

type World struct{
	rooms 	[]*Room
	lobby	*Room
	join	chan net.Conn
}

func (self *World) startLoop() {

	self.rooms = make([]*Room, 0)
	self.lobby = NewRoom(5)

	var joinReq JoinRequest
	pendingJoinReqs := make([]JoinRequest, 0)

	for {

		// Here, if we have any pending requests, we retrieve the current lobby's join
		// channel. Whether the request is actually sent to the channel or not depends
		// entirely on whether the lobby is still open. If the current lobby has been
		// closed at this point, `currentLobbyJoin` will be nil until we recieve from
		// its `Started` channel and have invalidated the lobby.

		var currentLobbyJoin chan<- JoinRequest

		if len(pendingJoinReqs) > 0 {
			joinReq = pendingJoinReqs[0]
			currentLobbyJoin = self.lobby.Join()
		}

		select {

			case <-self.lobby.Started() :

				self.rooms = append(self.rooms, self.lobby)
				self.lobby = NewRoom(5)
				break

			case currentLobbyJoin <-joinReq :

				// Note that if the Lobby is closed, its join chan will be nil, turning off
				// this select case and guarding against any join requests being made until
				// its Started message has been processed and the lobby has been invalidated.

				pendingJoinReqs = pendingJoinReqs[1:]
				break

			case conn := <-self.join :

				pendingJoinReqs = append(pendingJoinReqs, JoinRequest{ Conn: conn, })
				break

		}
	}
}

func (self *World) Start() {
	go self.startLoop()
}

func (self *World) Register(conn net.Conn) {
	self.join <-conn
}

func NewWorld() *World {
	world := &World{
		join: make(chan net.Conn),
	}
	return world
}
