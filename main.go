package main

import (
	"flag"
	stdlog "log"
	"os"
	"os/signal"

	"github.com/wenruo95/gossip/biz/config"
	"github.com/wenruo95/gossip/biz/control"
	"github.com/wenruo95/gossip/pkg/log"
)

func main() {
	var conf string
	flag.StringVar(&conf, "conf", "./conf/common.yaml", "-conf=./conf/common.yaml")

	if err := config.InitConfig(conf); err != nil {
		stdlog.Fatalf("init config error:%v", err)
	}
	if err := log.InitLogger(config.LogConfig()); err != nil {
		stdlog.Fatalf("init log error:%v", err)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Kill, os.Interrupt)
		sig := <-ch
		control.Close(sig.String())
	}()

	if err := control.Run(); err != nil {
		log.Fatalf("error:%v", err)
	}
}
