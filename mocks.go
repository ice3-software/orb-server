
package main

import (
	"net"
	"time"
	"io"
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
	Written		chan bool
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
		if self.Written != nil {
			self.Written <-true
		}
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
