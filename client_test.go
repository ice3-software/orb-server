package main

import (
	"testing"
	"io"
	//"sync"
	"bytes"
)

func TestNewClientOrb(t *testing.T) {

	mockConn := NewMockConn(nil)
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

	bsonOrb := BSONFromOrb(t, Orb{
		X: 1.0,
		Y: 2.0,
		ID: "4aa4ec0e-cf3f-4f4d-827f-f1415b51d26d",
	})

	mockConn := NewMockConn(nil)
	mockConn.SetMockReadData(bsonOrb)

	client := NewOrbClient(mockConn, make(chan Orb))
	readOrb := <-client.Read()

	AssertOrbAt(t, readOrb, 1.0, 2.0)

	if readOrb.ID == "4aa4ec0e-cf3f-4f4d-827f-f1415b51d26d" {
		t.Errorf("ID should not be overwritten, got %q instead of original value.", readOrb.ID)
	}

}

func TestReadEOF(t *testing.T) {

	mockConn := NewMockConn(nil)
	mockConn.SetMockReadError(io.EOF)

	client := NewOrbClient(mockConn, make(chan Orb))

	if !<-client.Disconnect() {
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

	mockConn := NewMockConn(nil)
	client := NewOrbClient(mockConn, make(chan Orb))
	err := client.Close()

	if !mockConn.Closed() {
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
	bsonOrb := BSONFromOrb(t, testOrb)

	conn := NewMockConn(nil)
	conn.SetMockReadData(bsonOrb)

	client := NewOrbClient(conn, make(chan Orb))
	client.Write() <-testOrb

	// @note No guarentee that the write has actually finished here. We need
	// to block on a `written` channel in the connection and _then_ perform
	// the assertion

	if !bytes.Equal(conn.MockWriteData(), bsonOrb) {
		t.Errorf("Orb not marshalled correctly, got %q", conn.MockWriteData())
	}

}

func TestWriteEOF(t *testing.T) {

	conn := NewMockConn(nil)
	conn.SetMockWriteError(io.EOF)

	client := NewOrbClient(conn, make(chan Orb))
	client.Write() <-Orb{}

	if !<-client.Disconnect() {
		t.Errorf("Should have recieved true value")
	}

}
