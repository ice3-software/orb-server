package main

import (
	"testing"
	"net"
	"time"
	"io"
	//"sync"
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"errors"
)

//
// A mock implementation of net.Conn
// TODO: Mutex locks here. This mock is not thread safe
//
type MockConn struct {
	ReadData	[]byte
	ReadError	error
	WriteData 	[]byte
	WriteError 	error
	Closed		bool
	//lock		sync.RWMutex
}

func (self *MockConn) Read(b []byte) (n int, err error) {

	if self.ReadError != nil {
		err = self.ReadError
	} else {
		n = len(self.ReadData)
		copy(b, self.ReadData)
	}
	self.ReadError = io.EOF
	return
}

func (self *MockConn) Write(b []byte) (n int, err error) {

	if self.WriteError != nil {
		err = self.WriteError
	} else {
		n = len(b)
		self.WriteData = make([]byte, n)
		copy(b, self.WriteData)
	}
	self.WriteError = io.EOF
	return
}

func (self *MockConn) Close() error {

	self.Closed = true
	return errors.New("error")
}

func (self *MockConn) LocalAddr() net.Addr { return nil }
func (self *MockConn) RemoteAddr() net.Addr { return nil }
func (self *MockConn) SetDeadline(t time.Time) error { return nil }
func (self *MockConn) SetReadDeadline(t time.Time) error { return nil }
func (self *MockConn) SetWriteDeadline(t time.Time) error { return nil }

//
// OrbClient tests
//

func TestNewClientOrb(t *testing.T) {

	mockConn := &MockConn{}
	client := NewOrbClient(mockConn);

	if client == nil {
		t.Errorf("Client is nil")
	}
	if client.Orb().ID == "" {
		t.Errorf("Client's Orb ID must be UUID")
	}
	if client.Conn() != mockConn {
		t.Errorf("Client's connection was not injected")
	}
	if client.Write() == nil {
		t.Errorf("Client's Write channel is nil")
	}
	if client.Read() == nil {
		t.Errorf("Client's Read channel is nil")
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
	})
	generatedID := client.Orb().ID
	readOrb := <-client.Read()

	if readOrb.X != 1 {
		t.Errorf("X not unmarshaled correctly, got %f", readOrb.X)
	}
	if readOrb.Y != 2 {
		t.Errorf("Y not unmarshaled correctly, got %f", readOrb.Y)
	}
	if readOrb.ID != generatedID {
		t.Errorf("ID should not be overwritten, got %q instead of %q", readOrb.ID, generatedID)
	}

}

func TestReadEOF(t *testing.T) {

	client := NewOrbClient(&MockConn{
		ReadError: io.EOF,
	})
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
	client := NewOrbClient(conn)
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
	client := NewOrbClient(conn)
	client.Write() <-testOrb

	if bytes.Equal(conn.WriteData, bsonOrb) {
		t.Errorf("Orb not marshalled correctly, got %q", conn.WriteData)
	}

}

func TestWriteEOF(t *testing.T) {

	client := NewOrbClient(&MockConn{
		WriteError: io.EOF,
	})
	client.Write() <-Orb{}
	disconnect := <-client.Disconnect()

	if !disconnect {
		t.Errorf("Should have recieved true value")
	}

}
