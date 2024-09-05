package impl

import (
	// "encoding/gob"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/sshx/pkg/types"
)

type ProxyService struct {
	BaseImpl
	RemotePort int32
}

func (s *ProxyService) Code() int32 {
	return types.APP_TYPE_PROXY_SERVICE
}

func (s *ProxyService) Preper() error {
	return nil
}

func (s *ProxyService) Dial() error {
	return nil
}

func (s *ProxyService) Response() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	logrus.Debug("Response impl proxy service")

	logrus.Debug("Dial local addr ", s.RemotePort)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", s.RemotePort))
	if err != nil {
		return err
	}
	s.BaseImpl.conn = &conn
	return nil
}