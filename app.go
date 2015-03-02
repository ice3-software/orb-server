
package main

import (
	"net"
	"fmt"
	"os"
)


type OrbServer struct {
	Host		net.IP
	Port		uint16
}


func (self *OrbServer) BasePath() string {
	return fmt.Sprintf("%s:%d", self.Host.String(), self.Port)
}

func (self *OrbServer) Serve() error {

	fmt.Println("Listening to ", self.BasePath())
	server, err := net.Listen("tcp", self.BasePath())

	if err != nil {
		return err
	}

	for {

		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting incoming connection: ", err)
			continue
		}

		conn.Write([]byte("Welcome to orb!\n"))

		msgBuff := make([]byte, 1024)
		_, readErr := conn.Read(msgBuff)

		if readErr == nil {
			conn.Write(msgBuff)
		} else {
			conn.Write([]byte(readErr.Error()))
		}

		conn.Close()

	}

	defer server.Close()
	return nil
}

func NewOrbServer() *OrbServer {
	return &OrbServer{
		Host: net.IPv4(127, 0, 0, 1),
		Port: 9090,
	}
}

func main() {

	server := NewOrbServer()
	err := server.Serve()

	if err != nil {
		fmt.Println("Could not serve")
		os.Exit(1)
	}

}
