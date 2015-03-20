
package main

import (
	"net"
)

type World struct{
	rooms 	[]*Room
	lobby	*Room
	join	chan net.Conn
}

func (self *World) bumpLobby() {

	if self.lobby != nil {
		self.rooms = append(self.rooms, self.lobby)
	}
	self.lobby = NewRoom(1)
}

func (self *World) startLoop() {

	self.bumpLobby()
	self.rooms = make([]*Room, 0)

	var joinReq JoinRequest
	pendingJoinReqs := make([]JoinRequest, 0)

	for {

		// Here, if we have any pending requests, we retrieve the current lobby's join
		// channel. Whether the request is actually sent to the channel or not depends
		// entirely on whether the lobby is still open. If the current lobby has been
		// closed at this point, `currentLobbyJoin` will never be received upon because
		// the Room's internal join channel will have been nilled out. We need to then
		// recieve from its `Started` channel and invalidate the lobby.

		var currentLobbyJoin chan<- JoinRequest

		if len(pendingJoinReqs) > 0 {
			joinReq = pendingJoinReqs[0]
			currentLobbyJoin = self.lobby.Join()
		}

		select {

			case <-self.lobby.Started() :
				self.bumpLobby()
				break

			case currentLobbyJoin <-joinReq :

				// Note that if the Lobby is closed, its internal join chan will be nil, turning off
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
