package main

import (
	"testing"
	"net"
	"time"
	"io"
	//"fmt"
	"gopkg.in/mgo.v2/bson"
	//"bytes"
	//"errors"
)

//
// A mock implementation of net.Conn
//
type MockConn struct {
	ReadData	[]byte
	ReadError	error
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

func (self *MockConn) Write(b []byte) (n int, err error) { return }
func (self *MockConn) Close() error { return nil }
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

func TestReadOrbMessage(t *testing.T) {

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

	mockConn := &MockConn{
		ReadData: bsonBuf,
	}
	client := NewOrbClient(mockConn)
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

}
