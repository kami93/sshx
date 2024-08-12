package impl

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/sshx/pkg/conf"
	"github.com/suutaku/sshx/pkg/types"
)

// Request struct which send to Local TCP listenner
type Sender struct {
	Type       int32 // Request type defined on types
	PairId     []byte
	Detach     bool
	LocalEntry string
	Payload    []byte // Application specify payload
	Status     int32
	ProxyHostPort int32 // Host port for proxy app.       
}

func NewProxySender(imp *Proxy, optCode int32) *Sender {
	if imp.HostId() == "" {
		logrus.Warn("Host Id not set, maybe you should set it on Preper")
	}

	ret := &Sender{
		Type: (imp.Code() << flagLen) | optCode,
		ProxyHostPort: (imp.Hostport() << flagLen) | optCode,
	}

	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(imp)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	ret.Payload = buf.Bytes()
	cm := conf.NewConfManager("")
	ret.LocalEntry = fmt.Sprintf("127.0.0.1:%d", cm.Conf.LocalTCPPort)
	ret.PairId = []byte(imp.PairId())
	return ret
}

func NewSender(imp Impl, optCode int32) *Sender {
	if imp.HostId() == "" {
		logrus.Warn("Host Id not set, maybe you should set it on Preper")
	}

	ret := &Sender{
		Type: (imp.Code() << flagLen) | optCode,
	}

	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(imp)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	ret.Payload = buf.Bytes()
	cm := conf.NewConfManager("")
	ret.LocalEntry = fmt.Sprintf("127.0.0.1:%d", cm.Conf.LocalTCPPort)
	ret.PairId = []byte(imp.PairId())
	return ret
}

func (sender *Sender) GetAppCode() int32 {
	return sender.Type >> flagLen
}

func (sender *Sender) GetOptionCode() int32 {
	return sender.Type & 0xff
}

func (sender *Sender) GetImpl() Impl {

	appcode := sender.GetAppCode()
	impl := GetImpl(appcode)

	// if appcode == types.APP_TYPE_PROXY {
	// 	// hostport := (sender.ProxyHostPort >> flagLen)
	// 	var hostport int32 = 22
	// 	impl.SetHostport(hostport)
	// 	logrus.Warn("Setting host port for proxy ", hostport)
	// }

	buf := bytes.NewBuffer(sender.Payload)
	err := gob.NewDecoder(buf).Decode(impl)
	if err != nil {
		logrus.Error(err)
	}
	return impl
}

func (sender *Sender) Send() (net.Conn, error) {
	conn, err := net.Dial("tcp", sender.LocalEntry)
	if err != nil {
		return nil, err
	}
	err = gob.NewEncoder(conn).Encode(sender)
	if err != nil {
		return nil, err
	}
	logrus.Debug("waiting TCP Responnse")

	err = gob.NewDecoder(conn).Decode(sender)
	if err != nil {
		return nil, err
	}

	logrus.Debug("TCP Responnse OK ", string(sender.PairId))
	if sender.Status != 0 {
		return nil, fmt.Errorf("response error")
	}
	return conn, nil
}

func (sender *Sender) SendDetach() (net.Conn, error) {
	sender.Detach = true
	return sender.Send()

}
