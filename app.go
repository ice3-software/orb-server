
package main

import (
	"os"
	"fmt"
)

func main() {

	server := NewOrbServer()
	err := server.Listen()

	if err != nil {
		fmt.Println("Could not serve")
		os.Exit(1)
	}

}
