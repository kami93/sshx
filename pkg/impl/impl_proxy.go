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
	RemotePort int32
	Running     bool
	ProxyHostId string
}

func NewProxy(port int32, remoteport int32, host string) *Proxy {
	return &Proxy{
		ProxyPort:   port,
		RemotePort: remoteport,
		ProxyHostId: host,
	}
}

func (base *Proxy) Preper() error {
	logrus.Debug("Preper impl proxy")
	return nil
}

func (p *Proxy) Code() int32 {
	return types.APP_TYPE_PROXY
}

func (p *Proxy) Start() error {
	conf.ClearKnownHosts(fmt.Sprintf("127.0.0.1:%d", p.ProxyPort))
	p.Running = true
	listenner, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p.ProxyPort))
	if err != nil {
		return err
	}
	fmt.Println("Proxy for", p.ProxyHostId, ":", p.RemotePort, " at :", p.ProxyPort)

	for p.Running {
		conn, err := listenner.Accept()
		if err != nil {
			continue
		}
		// proxy.conn = &conn
		go p.doDial(conn)

	}
	logrus.Debug("Close proxy for ", p.ProxyHostId)

	return nil
}

func (p *Proxy) Response() error {
	logrus.Debug("Response impl proxy")
	return nil
}

func (p *Proxy) Close() {
	p.Running = false
	logrus.Debug("close proxy impl")
}

func (p *Proxy) GetRemotePort() int32 {
	return p.RemotePort
}

func (p *Proxy) SetRemotePort(port int32) error {
	p.RemotePort = port

	return nil
}

func (p *Proxy) doDial(inconn net.Conn) {
	imp := &ProxyService{
		BaseImpl: BaseImpl{
			HId:        p.ProxyHostId,
			ConnectNow: true,
		},
		RemotePort: p.RemotePort,
	}
	logrus.Debug("Dial to ", p.ProxyHostId, ":", p.RemotePort)

	imp.SetParentId(p.PairId())
	sender := NewSender(imp, types.OPTION_TYPE_UP)
	conn, err := sender.Send()
	logrus.Debug(err)

	if err != nil {
		logrus.Error(err)
		return
	}
	defer conn.Close()
	// defer func() {
	// 	conn.Close()
	// 	closeSender := NewSender(imp, types.OPTION_TYPE_DOWN)
	// 	closeSender.PairId = sender.PairId
	// 	closeSender.SendDetach()
	// }()
	utils.Pipe(&inconn, &conn)
}
