package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
	"github.com/suutaku/sshx/pkg/impl"
	"github.com/suutaku/sshx/pkg/types"
)

func cmdStopProxy(cmd *cli.Cmd) {
	cmd.Spec = "PID"
	pairId := cmd.StringArg("PID", "", "Connection pair id which can found by using status command")
	cmd.Action = func() {
		imp := impl.NewProxy(0, "", 0)
		imp.NoNeedConnect()
		sender := impl.NewSender(imp, types.OPTION_TYPE_DOWN)
		sender.PairId = []byte(*pairId)
		sender.SendDetach()
	}
}

func cmdStartProxy(cmd *cli.Cmd) {
	// cmd.Spec = "-P [-d] ADDR"
	cmd.Spec = "-L -R ADDR"
	localPort := cmd.IntOpt("L", 0, "local proxy port")
	remotePort := cmd.IntOpt("R", 0,    "remote proxy port")

	addr := cmd.StringArg("ADDR", "", "remote target address [username]@[host]:[port]")
	cmd.Action = func() {
		if localPort == nil || *localPort == 0 {
			fmt.Println("please set a local port")
		}
		if remotePort == nil || *remotePort == 0 {
			fmt.Println("please set a remote port")
		}
		if addr == nil || *addr == "" {
			fmt.Println("please set a remote device")
		}

		proxy := impl.NewProxy(int32(*localPort), *addr, int32(*remotePort))
		proxy.Preper()
		// proxy.NoNeedConnect()

		fmt.Println("pass 1")

		sender := impl.NewSender(proxy, types.OPTION_TYPE_UP)

		fmt.Println("pass 2")

		// _, err := sender.SendDetach()
		conn, err := sender.Send()
		
		fmt.Println("pass 3")

		if err != nil {
			logrus.Error(err)
		}
		err = proxy.Start(conn)
		
		fmt.Println("pass 4")

		if err != nil {
			logrus.Error(err)
		}
		proxy.Close()
	}
}

func cmdProxy(cmd *cli.Cmd) {
	cmd.Command("start", "start proxy service", cmdStartProxy)
	cmd.Command("stop", "stop proxy service", cmdStopProxy)
}
