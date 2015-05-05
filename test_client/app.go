
package intest

import (
	"net/http"
	"time"
)

const (
	NumClientConns = 300,
	NumMessagesPerSec = 40,
)

const (
	ServerAddress = "127.0.0.1:10101"
	ServerProtocol = "tcp"
)

func readLoop(exit chan bool, read chan string, conn net.Conn) {

	for {

	}
}

func writeLoop(exit chan bool, write chan string, conn net.Conn) {

	wait := time.Wait(time.Second/NumMessagesPerSec)

	for {
		<-wait

	}

}

func createClient(exit chan bool) error {

	if conn, err := net.Dial(ServerProtocol, ServerAddress); err != nil {
		return err
	}

	readExit := make(chan bool)
	writeExit := make(chan bool)

	go func(exit chan bool) {
		for exitNot := range exit {
			readExit <-exitNot
			writeExit <-exitNot
		}
	}(exit)

	go readLoop(readExit, conn)
	go writeLoop(writeExit, conn)

}

func main() {



}
