package utils

import (
	"github.com/wenruo95/gossip/pkg/log"
)

type ExecFunc func() error

type ExecScope interface {
	With() ExecScope
	Exec() error
}

type execItem struct {
	key    string
	fn     ExecFunc
	sync   bool
	ignore bool
}
type execChain struct {
	execList []*execItem
}

func NewExecChain() *execChain {
	chain := new(execChain)
	chain.execList = make([]*execItem, 0)
	return chain
}

func (chain *execChain) With(key string, fn ExecFunc) *execChain {
	chain.execList = append(chain.execList, &execItem{key: key, fn: fn})
	return chain
}

func (chain *execChain) WithGo(key string, fn ExecFunc) *execChain {
	chain.execList = append(chain.execList, &execItem{key: key, fn: fn, sync: true})
	return chain
}
func (chain *execChain) WithIgnore(key string, fn ExecFunc, ignoreErr bool) *execChain {
	chain.execList = append(chain.execList, &execItem{key: key, fn: fn, ignore: true})
	return chain
}

func (chain *execChain) Exec() error {
	for idx, item := range chain.execList {
		if item.sync {
			go func() {
				if err := item.fn(); err != nil {
					log.Warnf("exec idx:%v key:%v error:%v", idx, item.key, err)
				}
			}()
			continue
		}

		if err := item.fn(); err != nil {
			log.Warnf("exec idx:%v key:%v error:%v", idx, item.key, err)
			if item.ignore {
				continue
			}
			return err
		}
	}
	return nil
}
