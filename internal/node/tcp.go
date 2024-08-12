package node

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/sshx/pkg/impl"
	"github.com/suutaku/sshx/pkg/types"
)

func (node *Node) ServeTCP() {
	listenner, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", node.confManager.Conf.LocalTCPPort))
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	defer listenner.Close()
	for node.running {
		sock, err := listenner.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		tmp := impl.Sender{}
		err = gob.NewDecoder(sock).Decode(&tmp)
		if err != nil {
			logrus.Warn("read not ok", err)
			sock.Close()
			continue
		}
		switch tmp.GetOptionCode() {
		case types.OPTION_TYPE_UP:
			logrus.Warn("up option")
			impl := tmp.GetImpl()
			if impl == nil {
				logrus.Error("unkwon implementation")
				continue
			}
			poolId := types.NewPoolId(time.Now().UnixNano(), impl.Code())
			err := node.connMgr.CreateConnection(&tmp, sock, *poolId)
			if err != nil {
				sock.Close()
				logrus.Error(err)
			}

		case types.OPTION_TYPE_DOWN:
			logrus.Warn("down option ", string(tmp.PairId))
			err := node.connMgr.DestroyConnection(&tmp, sock)
			if err != nil {
				logrus.Error(err)
			}

		case types.OPTION_TYPE_STAT:
			logrus.Warn("stat option")
			err := node.connMgr.Status(tmp, sock)
			if err != nil {
				sock.Close()
				logrus.Error(err)
			}
		case types.OPTION_TYPE_ATTACH:
			logrus.Warn("attach option")
			err := node.connMgr.AttachConnection(&tmp, sock)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
