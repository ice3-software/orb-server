package main

import (
	"testing"
	"io"
	//"sync"
	"gopkg.in/mgo.v2/bson"
	"bytes"
)

func TestNewClientOrb(t *testing.T) {

	mockConn := &MockConn{}
	read := make(chan Orb)
	client := NewOrbClient(mockConn, read);

	if client == nil {
		t.Errorf("Client is nil")
	}
	if client.Conn() != mockConn {
		t.Errorf("Client's connection was not injected")
	}
	if client.Write() == nil {
		t.Errorf("Client's Write channel is nil")
	}
	if client.Read() != read {
		t.Errorf("Client's Read channel was not injected")
	}
	if client.Disconnect() == nil {
		t.Errorf("Client's Disconnect channel is nil")
	}

}

func TestReadOrb(t *testing.T) {

	bsonBuf, err := bson.Marshal(&Orb{
		X: 1.0,
		Y: 2.0,
		ID: "4aa4ec0e-cf3f-4f4d-827f-f1415b51d26d",
	})

	if err != nil {
		t.Error("Error marshalling the test data: ", err)
	} else {
		t.Log("Marshalled test data: ", bsonBuf)
	}

	client := NewOrbClient(&MockConn{
		ReadData: bsonBuf,
	}, make(chan Orb))
	readOrb := <-client.Read()

	if readOrb.X != 1 {
		t.Errorf("X not unmarshaled correctly, got %f", readOrb.X)
	}
	if readOrb.Y != 2 {
		t.Errorf("Y not unmarshaled correctly, got %f", readOrb.Y)
	}
	if readOrb.ID == "4aa4ec0e-cf3f-4f4d-827f-f1415b51d26d" {
		t.Errorf("ID should not be overwritten, got %q instead of original value.", readOrb.ID)
	}

}

func TestReadEOF(t *testing.T) {

	client := NewOrbClient(&MockConn{
		ReadError: io.EOF,
	}, make(chan Orb))
	disconnect := <-client.Disconnect()

	if !disconnect {
		t.Errorf("Should have recieved true value")
	}

}

func TestReadError(t *testing.T) {
	//t.Fail()
}

func TestReadUnmarshalError(t *testing.T) {
	//t.Fail()
}

func TestClose(t *testing.T) {

	conn := &MockConn{}
	client := NewOrbClient(conn, make(chan Orb))
	err := client.Close()

	if !conn.Closed {
		t.Errorf("Should have closed connection")
	}
	if err == nil {
		t.Errorf("Should have returned error from connection")
	}

}

func TestWriteOrb(t *testing.T) {

	testOrb := Orb{
		X: 5.0,
		Y: 6.0,
		ID: "4aa4ec0e-cf3f-4f4d-827f-f1415b51d26d",
	}
	bsonOrb, err := bson.Marshal(&testOrb)

	if err != nil {
		t.Error("Error marshalling the test data: ", err)
	} else {
		t.Log("Marshalled test data: ", bsonOrb)
	}

	conn := &MockConn{}
	client := NewOrbClient(conn, make(chan Orb))
	client.Write() <-testOrb

	if bytes.Equal(conn.WriteData, bsonOrb) {
		t.Errorf("Orb not marshalled correctly, got %q", conn.WriteData)
	}

}

func TestWriteEOF(t *testing.T) {

	client := NewOrbClient(&MockConn{
		WriteError: io.EOF,
	}, make(chan Orb))
	client.Write() <-Orb{}
	disconnect := <-client.Disconnect()

	if !disconnect {
		t.Errorf("Should have recieved true value")
	}

}
