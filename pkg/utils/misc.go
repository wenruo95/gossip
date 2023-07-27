package utils

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wenruo95/gossip/pkg/log"
)

func SignalServe(wait time.Duration, stop func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Info("recv signal:" + sig.String())

	go stop()
	if wait.Seconds() > 10 {
		wait = 10 * time.Second
	}
	time.Sleep(wait)
	log.Info("wait for " + wait.String() + " close")
	os.Exit(128 + int(sig.(syscall.Signal)))
}

var randstr = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandCharacter(n int) []byte {
	data := make([]byte, n)
	for i := 0; i < n; i++ {
		data[i] = randstr[rand.Int()%len(randstr)]
	}
	return data
}
