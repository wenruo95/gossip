package utils

import (
	"fmt"
)

type ExecFunc func() error

type ExecScope interface {
	With() ExecScope
	Exec() error
}

type execItem struct {
	key string
	fn  ExecFunc
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

func (chain *execChain) Exec() error {
	for idx, item := range chain.execList {
		if err := item.fn(); err != nil {
			return fmt.Errorf("exec idx:%v key:%v error:%w", idx, item.key, err)
		}
	}
	return nil
}
