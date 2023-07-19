package config

import "github.com/wenruo95/gossip/pkg/log"

var ConfPath string

func InitConfig() error {

	return nil
}

type ServerConfig struct {
	Addr      string
	TimeoutMs int64
}

func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr:      ":5298",
		TimeoutMs: 10,
	}
}

func GetLogConfig() *log.Config {
	cfg := new(log.Config)
	cfg.FileName = ""
	return cfg
}
