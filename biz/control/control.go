package control

import (
	"errors"
	"os"
	"os/signal"
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
	return utils.NewExecChain().
		WithGo("ticker", ctrl.ticker.Serve).
		WithGo("signal", signalServe).
		With("server", ctrl.server.Serve).
		Exec()
}

type control struct {
	server *tcp.Server
	ticker *ticker
	pcache *peerCache
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

func signalServe() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Kill, os.Interrupt)
	sig := <-ch
	log.Info("recv signal:" + sig.String())
	ctrl.Close()
	return nil
}

func (ctrl *control) Close() error {
	if ctrl == nil || ctrl.server == nil {
		return nil
	}
	return ctrl.server.Close()
}
