
package main

import (
	"net"
	"time"
	"io"
	"errors"
	"sync"
)

//
// A mock implementation of net.Conn
//
type MockConn struct {

	lock			sync.RWMutex

	mockReadData	[]byte
	mockReadError	error

	mockWriteData 	[]byte
	mockWriteError 	error

	closed			bool
	closedErr		error

	written			chan bool

}


// @note - Does not copy the errors so assumed that they are either immutable or are
// thread safe.

func (self *MockConn) SetMockWriteError(err error) {

	self.lock.Lock()
	defer self.lock.Unlock()

	self.mockWriteError = err
}

func (self *MockConn) SetMockReadError(err error) {

	self.lock.Lock()
	defer self.lock.Unlock()

	self.mockWriteError = err
}

func (self *MockConn) SetMockReadData(mockReadData []byte) {

	self.lock.Lock()
	defer self.lock.Unlock()

	self.mockReadData = make([]byte, len(mockReadData))
	copy(self.mockReadData, mockReadData)
}

func (self *MockConn) MockWriteData() []byte {

	self.lock.RLock()
	defer self.lock.RUnlock()

	mockWriteData := make([]byte, len(self.mockWriteData))
	copy(mockWriteData, self.mockWriteData)

	return mockWriteData
}

func (self *MockConn) Read(b []byte) (n int, err error) {

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.mockReadError != nil {
		err = self.mockReadError
	} else {
		n = len(self.mockReadData)
		copy(b, self.mockReadData)
	}
	self.mockReadError = io.EOF

	return
}

func (self *MockConn) Write(b []byte) (n int, err error) {

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.mockWriteError != nil {
		err = self.mockWriteError
	} else {
		n = len(b)
		self.mockWriteData = make([]byte, n)
		copy(self.mockWriteData, b)
		self.broadcastWrite(n)
	}

	return
}

func (self *MockConn) broadcastWrite(n int) {

	// @note - The problem with this is that it _assumes_ we're running not on the
	// recieving thread, thus if Write is executed on the same thread that recieves
	// from this chan we get deadlock.

	if self.written != nil && n > 0 {
		self.written <-true
	}
}

func (self *MockConn) Close() error {

	self.lock.Lock()
	defer self.lock.Unlock()

	self.closed = true
	return self.closedErr
}

func (self *MockConn) Closed() bool {

	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.closed
}

func (self *MockConn) LocalAddr() net.Addr { return nil }
func (self *MockConn) RemoteAddr() net.Addr { return nil }
func (self *MockConn) SetDeadline(t time.Time) error { return nil }
func (self *MockConn) SetReadDeadline(t time.Time) error { return nil }
func (self *MockConn) SetWriteDeadline(t time.Time) error { return nil }

func NewMockConn(written chan bool) *MockConn {
	return &MockConn{
		mockWriteData: 	make([]byte, 2048),
		closedErr: 		errors.New("Error on close"),
		written:		written,
	}
}
