package main

import (
	"flag"

	"github.com/wenruo95/gossip/biz/config"
	"github.com/wenruo95/gossip/biz/control"
	"github.com/wenruo95/gossip/pkg/log"
)

func init() {
	flag.StringVar(&config.ConfPath, "conf", "./conf/common.yaml", "-conf=./conf/common.yaml")
	flag.StringVar(&config.ServerAddr, "addr", "", "-addr=:5298")
	flag.Parse()
}

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("init config error:%v", err)
	}
	if err := log.InitLogger(config.GetLogConfig()); err != nil {
		log.Fatalf("init log error:%v", err)
	}
	defer log.Sync()

	log.Fatalf("error:%v", control.Run())
}
