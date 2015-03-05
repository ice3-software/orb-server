
package main

import (
	"testing"
	//"fmt"
	//"errors"
)

func TestNewRoom(t *testing.T) {

	room := NewRoom(1)

	if room == nil {
		t.Error("Room not created")
	}
}

func TestJoinRoom(t *testing.T) {

	room := NewRoom(1)
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

	bsonOrb := BSONFromOrb(t, Orb{
		X: 1.0,
		Y: 2.0,
		ID: "4aa4ec0e-cf3f-4f4d-827f-f1415b51d26d",
	})

	room := NewRoom(3)
	writtenChan := make(chan bool)

	reciever1 := NewMockConn(writtenChan)
	reciever2 := NewMockConn(writtenChan)
	sender := NewMockConn(writtenChan)
	sender.SetMockReadData(bsonOrb)

	room.Join() <-JoinRequest{ Conn: reciever1, }
	room.Join() <-JoinRequest{ Conn: reciever2, }
	room.Join() <-JoinRequest{ Conn: sender, }

	<-writtenChan
	<-writtenChan

	orb1 := OrbFromBSON(t, reciever1.MockWriteData())
	orb2 := OrbFromBSON(t, reciever2.MockWriteData())

	AssertOrbAt(t, orb1, 1.0, 2.0)
	AssertOrbAt(t, orb2, 1.0, 2.0)

}
