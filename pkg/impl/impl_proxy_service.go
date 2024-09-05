package impl

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/sshx/pkg/types"
)

type ProxyService struct {
	BaseImpl
	RemotePort int32
	isRuning   bool
	sendChan   chan PortMessage
	attachConn net.Conn
}

type PortMessage struct {
	Payload int32
}

func (s *ProxyService) Code() int32 {
	return types.APP_TYPE_PROXY_SERVICE
}

func (s *ProxyService) Preper() error {
	logrus.Debug("Preper impl proxy service")
	s.isRuning = true
	_c, _s := net.Pipe()
	s.attachConn = _c
	s.BaseImpl.conn = &_s

	go s.serveSend()

	return nil
}

func (s *ProxyService) Dial() error {
	return nil
}

func (s *ProxyService) serveRecv() {
	for {
		var msg PortMessage
		logrus.Debug("waiting message")
		err := gob.NewDecoder(s.attachConn).Decode(&msg)
		logrus.Debug("waiting message ok")
		if err != nil {
			logrus.Error(err)
			s.Close()
			return
		}
		// logrus.Debug("message come ", msg.Payload)
		logrus.Debug("set remote port", msg.Payload)
		s.RemotePort = msg.Payload
	}
}

func (s *ProxyService) serveSend() {
	if s.sendChan == nil {
		s.sendChan = make(chan PortMessage, s.RemotePort)
	}
	for {
		msg := <- s.sendChan
		err := gob.NewEncoder(s.attachConn).Encode(msg)
		if err != nil {
			logrus.Error(err)
			s.Close()
			return
		}
	}
}

func (s *ProxyService) Response() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	logrus.Debug("Response impl proxy service")

	// a naive test code
	s.isRuning = true
	_c, _s := net.Pipe()
	s.attachConn = _c
	s.BaseImpl.conn = &_s
	go s.serveRecv()

	logrus.Debug("Dial local addr ", s.RemotePort)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", s.RemotePort))
	if err != nil {
		return err
	}
	s.BaseImpl.conn = &conn
	return nil
}