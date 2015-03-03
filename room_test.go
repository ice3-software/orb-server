
package main

import (
	"testing"
)

func TestNewRoom(t *testing.T) {

	room := NewRoom()

	if room == nil {
		t.Error("Room not created")
	}
}

func TestJoinRoom(t *testing.T) {

	//room := NewRoom()

}
