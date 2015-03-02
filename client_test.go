package main

import (
	"testing"
	"net"
	"time"
)

//
// A mock implementation of net.Conn
//
type MockConn struct {}

func (self *MockConn) Read(b []byte) (n int, err error) { return }
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
