package utils

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecChain(t *testing.T) {

	var b1, b2, b3 bool
	chain := NewExecChain().
		With("test1", func() error {
			b1 = true
			return nil
		}).
		With("test2", func() error {
			b2 = true
			return errors.New("mock test2 error")
		}).
		With("test3", func() error {
			b3 = true
			return errors.New("mock test2 error")
		})
	err := chain.Exec()
	assert.NotNil(t, err)
	assert.True(t, b1)
	assert.True(t, b2)
	assert.False(t, b3)
	if err != nil {
		assert.True(t, strings.Contains(err.Error(), "mock test2"))
	}

}
