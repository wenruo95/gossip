package control

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/wenruo95/gossip/biz/config"
	"github.com/wenruo95/gossip/pkg/tcp"
	"github.com/wenruo95/gossip/pkg/utils"
)

var ctrl *control
var once int32

func Run() error {
	if swap := atomic.CompareAndSwapInt32(&once, 0, 1); !swap {
		return errors.New("error:dumplicate run control")
	}

	return utils.NewExecChain().
		With("control", initControl).
		WithGo("ticker", tickerServe).
		With("server", serverServe).
		Exec()
}

func Close() error {
	return ctrl.Close()
}

type control struct {
	server *tcp.Server
	ticker *ticker
}

func initControl() error {
	ctrl = new(control)

	cfg := config.GetServerConfig()
	ctrl.server = tcp.NewServer(
		tcp.WithAddr(cfg.Addr),
		tcp.WithHandler(ctrl),
		tcp.WithTimeout(time.Duration(cfg.TimeoutMs)*time.Millisecond),
	)
	ctrl.ticker = newTicker()
	return nil
}

func serverServe() error {
	return ctrl.server.Serve()
}

func tickerServe() error {
	return ctrl.ticker.Serve()
}

func (ctrl *control) Close() error {
	if ctrl == nil || ctrl.server == nil {
		return nil
	}
	return ctrl.server.Close()
}
