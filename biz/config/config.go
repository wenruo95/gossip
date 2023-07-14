package config

import "github.com/wenruo95/gossip/pkg/log"

func InitConfig(path string) error {

	return nil
}

func LogConfig() *log.Config {
	cfg := new(log.Config)
	cfg.FileName = "./log/gossip.log"
	return cfg
}
