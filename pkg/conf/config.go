package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/suutaku/sshx/internal/utils"
)

type Configure struct {
	LocalSSHPort        int32
	LocalHTTPPort       int32
	LocalTCPPort        int32
	ID                  string
	SignalingServerAddr string
	RTCConf             webrtc.Configuration
	ETHAddr             string
}

type ConfManager struct {
	Conf  *Configure
	Viper *viper.Viper
	Path  string
}

var defaultConfig = Configure{
	LocalHTTPPort:       80,
	LocalSSHPort:        22,
	LocalTCPPort:        2224,
	ID:                  uuid.New().String(),
	SignalingServerAddr: "http://alindev.kaist.ac.kr:5003",
	RTCConf: webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun1.l.google.com:19302",
					"stun:stun2.l.google.com:19302",
					"stun:stun3.l.google.com:19302",
					"stun:stun4.l.google.com:19302",
				},
			},
		},
	},
}

func ClearKnownHosts(subStr string) {
	subStr = strings.Replace(subStr, "127.0.0.1", "[127.0.0.1]", 1)
	//[127.0.0.1]:2222
	fileName := os.Getenv("HOME") + "/.ssh/known_hosts"
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Error(err)
		return
	}
	lines := strings.Split(string(input), "\n")
	var newLines []string
	for i, line := range lines {
		if strings.Contains(line, subStr) {
		} else {
			newLines = append(newLines, lines[i])
		}
	}
	output := strings.Join(newLines, "\n")
	err = ioutil.WriteFile(fileName, []byte(output), 0777)
	if err != nil {
		logrus.Error(err)
		return
	}
	//ioutil.WriteFile(fileName, []byte(res), 544)
}

func NewConfManager(homePath string) *ConfManager {
	if homePath == "" {
		homePath = utils.GetSSHXHome()
	}
	var tmp Configure
	vp := viper.New()
	vp.SetConfigName(".sshx_config")
	vp.SetConfigType("json")
	vp.AddConfigPath(homePath)
	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		err := vp.Unmarshal(&tmp)
		if err != nil {
			logrus.Error(err)
			return
		}
	})
	err := vp.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			defaultConfig.RTCConf.PeerIdentity = utils.HashString(fmt.Sprintf("%s%d", defaultConfig.ID, time.Now().Unix()))
			bs, _ := json.MarshalIndent(defaultConfig, "", "  ")
			vp.ReadConfig(bytes.NewBuffer(bs))
			err = vp.WriteConfigAs(path.Join(homePath, "./.sshx_config.json"))
			if err != nil {
				logrus.Error(err)
				os.Exit(1)
			}
			os.Chmod(path.Join(homePath, "./.sshx_config.json"), 0777)
			os.Exit(1)
		} else {
			logrus.Error(err)
			os.Exit(1)
		}
	}

	err = vp.Unmarshal(&tmp)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	ClearKnownHosts(fmt.Sprintf("127.0.0.1:%d", tmp.LocalSSHPort))
	return &ConfManager{
		Conf:  &tmp,
		Viper: vp,
		Path:  homePath,
	}
}

func (cm *ConfManager) Set(key, value string) {
	logrus.Info("key/value", key, value)
	cm.Viper.Set(key, value)
	err := cm.Viper.Unmarshal(cm.Conf)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = cm.Viper.WriteConfig()
	if err != nil {
		logrus.Error(err)
		return
	}
}

func (cm *ConfManager) Show() {
	bs, _ := json.MarshalIndent(cm.Conf, "", "  ")
	logrus.Info("read configure file at: ", cm.Path+"/.sshx_config.json")
	logrus.Info(string(bs))
}
