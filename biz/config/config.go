package config

import "github.com/wenruo95/gossip/pkg/log"

var ConfPath string
var ServerAddr string

func InitConfig() error {

	return nil
}

type ControlConfig struct {
	CenterList []string
}

type ServerConfig struct {
	Addr      string
	TimeoutMs int64
}

func GetServerConfig() *ServerConfig {
	cfg := &ServerConfig{
		Addr:      ":5298",
		TimeoutMs: 10,
	}
	if len(ServerAddr) > 0 {
		cfg.Addr = ServerAddr
	}
	return cfg
}

func GetLogConfig() *log.Config {
	cfg := new(log.Config)
	cfg.FileName = ""
	return cfg
}
