package impl

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/sshx/internal/utils"
	"github.com/suutaku/sshx/pkg/conf"
	"github.com/suutaku/sshx/pkg/types"
)

type Proxy struct {
	BaseImpl
	ProxyPort   int32
	Running     bool
	ProxyHostId string
	ProxyHostPort int32
}

func NewProxy(port int32, host string, hostport int32) *Proxy {
	return &Proxy{
		ProxyPort:   port,
		ProxyHostId: host,
		ProxyHostPort: hostport
	}
}

func (p *Proxy) Code() int32 {
	return types.APP_TYPE_PROXY
}

func (p *Proxy) Start(conn *net.Conn) error {
	conf.ClearKnownHosts(fmt.Sprintf("127.0.0.1:%d", p.ProxyPort))
	p.Running = true
	listenner, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p.ProxyPort))
	if err != nil {
		return err
	}
	fmt.Println("Proxy for ", p.ProxyHostId, " at :", p.ProxyPort)

	for p.Running {
		inconn, err := listenner.Accept()
		if err != nil {
			continue
		}
		
		defer conn.Close()
		utils.Pipe(&inconn, &conn)
	}
	logrus.Debug("Close proxy for ", p.ProxyHostId)

	return nil
}

func (p *Proxy) Close() {
	p.Running = false
	logrus.Debug("close proxy impl")
}

func (s *Proxy) Response() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	
	logrus.Debug("Dial local addr ", s.ProxyHostPort)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", s.ProxyHostPort))
	if err != nil {
		return err
	}
	s.BaseImpl.conn = &conn
	return nil
}

// func (p *Proxy) doDial(inconn net.Conn) {
// 	imp := &SSH{
// 		BaseImpl: BaseImpl{
// 			HId:        p.ProxyHostId,
// 			ConnectNow: true,
// 		},
// 	}
// 	imp.SetParentId(p.PairId())
// 	sender := NewSender(imp, types.OPTION_TYPE_UP)
// 	conn, err := sender.Send()
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// 	defer conn.Close()
// 	// defer func() {
// 	// 	conn.Close()
// 	// 	closeSender := NewSender(imp, types.OPTION_TYPE_DOWN)
// 	// 	closeSender.PairId = sender.PairId
// 	// 	closeSender.SendDetach()
// 	// }()
// 	utils.Pipe(&inconn, &conn)
// }
