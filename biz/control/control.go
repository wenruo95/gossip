package control

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/wenruo95/gossip/biz/config"
	"github.com/wenruo95/gossip/pkg/log"
	"github.com/wenruo95/gossip/pkg/tcp"
	"github.com/wenruo95/gossip/pkg/utils"
)

var ctrl *control
var once int32

func Run() error {
	if swap := atomic.CompareAndSwapInt32(&once, 0, 1); !swap {
		return errors.New("error:dumplicate run control")
	}

	ctrl = newControl()
	defer ctrl.Close()

	return utils.NewExecChain().
		WithGo("ticker", ctrl.ticker.Serve).
		WithGo("server", ctrl.server.Serve).
		With("signal", signalServe).
		Exec()
}

type control struct {
	server *tcp.Server
	ticker *ticker
	pcache *peerCache
}

func signalServe() error {
	utils.SignalServe(2*time.Second, ctrl.Close)
	return nil
}

func newControl() *control {
	c := new(control)
	serverConfig := config.GetServerConfig()
	c.server = tcp.NewServer(
		tcp.WithHandler(c),
		tcp.WithAddr(serverConfig.Addr),
		tcp.WithTimeout(time.Duration(serverConfig.TimeoutMs)*time.Millisecond),
	)
	c.pcache = newPeerCache()
	c.ticker = newTicker()
	return c
}

func (ctrl *control) Close() {
	if ctrl == nil {
		return
	}
	if err := ctrl.server.Close(); err != nil {
		log.Error("close server error:" + err.Error())
	}
}
