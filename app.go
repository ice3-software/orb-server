
package main

import (
	"net"
	"fmt"
	"os"
)


type OrbClient struct {
	X			float32
	Y 			float32
	Conn 		*net.TCPConn
	Room		*Room
}


type Room struct {
	Clients []*OrbClient
}

func (self *Room) Join(client *OrbClient) {
	self.Clients = append(self.Clients, client)
	if self.Full() {
		self.notifyFull()
	}
}

func (self *Room) Count() int {
	return len(self.Clients)
}

func (self *Room) notifyFull() {
	for _, client := range self.Clients {
		client.Conn.Write([]byte("Let the games begin!"))
	}
}

func (self *Room) Full() bool {
	return len(self.Clients) > 5
}

func NewRoom() *Room {
	return &Room{
		Clients: make([]*OrbClient, 0),
	}
}


type World struct{
	Rooms []*Room
}

func (self *World) Count() int {
	return len(self.Rooms)
}

func (self *World) WaitingRoom() (room *Room) {

	if self.Count() == 0 {
		room = NewRoom()
		self.Rooms = append(self.Rooms, room)
		fmt.Println("First waiting room")
	} else {
		fmt.Println("Waiting room... ", self.Count())
		room = self.Rooms[self.Count() - 1]
	}

	return
}

func (self *World) Register(conn *net.TCPConn) *OrbClient {

	fmt.Println("Registering client ", self.Count())

	waitingRoom := self.WaitingRoom()
	client := &OrbClient{
		Conn: conn,
		Room: waitingRoom,
	}
	waitingRoom.Join(client)

	if waitingRoom.Full() {
		fmt.Println("Waiting room if full ! Notifying ", self.Count())
		self.Rooms = append(self.Rooms, NewRoom())
	}

	return client

}

func NewWorld() *World {
	return &World{
		Rooms: make([]*Room, 0),
	}
}


type OrbServer struct {

	Host		net.IP
	Port		uint16

	Conns		chan *net.TCPConn

	World		*World

}

func (self *OrbServer) BasePath() string {
	return fmt.Sprintf("%s:%d", self.Host.String(), self.Port)
}

func (self *OrbServer) Listen() error {

	fmt.Println("Listening to ", self.BasePath())
	server, err := net.Listen("tcp", self.BasePath())

	if err != nil {
		return err
	}

	go self.Serve()

	for {

		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting incoming connection: ", err)
		} else {
			self.Conns <- conn.(*net.TCPConn)
		}

	}

	defer server.Close()
	return nil
}

func (self *OrbServer) Serve() {

	self.World = NewWorld()

	for {

		select {
		case conn := <-self.Conns:

			fmt.Println("New connection")
			conn.Write([]byte("Welcome to orb! You have been registered in the world\n"))
			self.World.Register(conn)

			//msgBuff := make([]byte, 1024)
			//_, readErr := conn.Read(msgBuff)

			//if readErr == nil {
			//	conn.Write(msgBuff)
			//} else {
			//	conn.Write([]byte(readErr.Error()))
			//}

		break

		}
	}
}

func NewOrbServer() *OrbServer {
	return &OrbServer{
		Host: net.IPv4(127, 0, 0, 1),
		Port: 9090,
		Conns: make(chan *net.TCPConn),
	}
}

func main() {

	server := NewOrbServer()
	err := server.Listen()

	if err != nil {
		fmt.Println("Could not serve")
		os.Exit(1)
	}

}
