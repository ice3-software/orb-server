
package main

import (
	"testing"
	"bytes"
)

func TestNewRoom(t *testing.T) {

	room := NewRoom()

	if room == nil {
		t.Error("Room not created")
	}
}

func TestJoinRoom(t *testing.T) {

	room := NewRoom()
	mockConn := &MockConn{}
	joined := make(chan bool)

	room.Join() <-JoinRequest{
		Conn: mockConn,
		Joined: joined,
	}

	if !<-joined {
		t.Error("Expecting true")
	}

}

func TestReadsBroadcastToRoom(t *testing.T) {

	bsonOrb := BSONOrb(t)
	room := NewRoom()
	sender := &MockConn{
		ReadData: bsonOrb,
	}
	reciever1 := &MockConn{
		Written: make(chan bool),
	}
	reciever2 := &MockConn{
		Written: make(chan bool),
	}

	room.Join() <-JoinRequest{ Conn: sender, }
	room.Join() <-JoinRequest{ Conn: reciever1, }
	room.Join() <-JoinRequest{ Conn: reciever2, }

	<-reciever1.Written
	//<-reciever2.Written

	if !bytes.Equal(reciever1.WriteData, bsonOrb) {
		t.Errorf("Expected %q, got %q", bsonOrb, reciever1.WriteData)
	}
	//if !bytes.Equal(reciever2.WriteData, bsonOrb) {
	//	t.Errorf("Expected %q, got %q", bsonOrb, reciever2.WriteData)
	//}

}
