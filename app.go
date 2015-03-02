
package main

import (
	"net"
	"fmt"
	"os"
)


type OrbServer struct {
	Host		net.IP
	Port		uint16
	Conns		chan *net.TCPConn
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

	for {

		select {

		case conn := <-self.Conns:

			conn.Write([]byte("Welcome to orb!\n"))
			msgBuff := make([]byte, 1024)
			_, readErr := conn.Read(msgBuff)

			if readErr == nil {
				conn.Write(msgBuff)
			} else {
				conn.Write([]byte(readErr.Error()))
			}

			conn.Close()

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
