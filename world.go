
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
	self.lobby = nil

	for {
		select {

			case conn := <-self.join
			if self.lobby == nil {
				self.lobby = NewRoom()
			}
			self.lobby.Join() <-JoinRequest{ Conn: conn, }
			break

			case <-self.lobby.Started()
			self.rooms = append(self.rooms, self.lobby)
			self.lobby = nil
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
		join: make(chan net.Conn)
	}
	return world
}
