package utils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZip(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		//str := RandCharacter(1024 * 1024 * 10)
		str := RandCharacter(10)
		buff, err := Zip([]byte(str))
		if err != nil {
			t.Errorf("zip error:%v", err)
			return
		}
		//t.Logf("str:%v buff:%v", str, string(buff))

		result, err := Unzip(buff)
		if err != nil {
			t.Errorf("unzip error:%v", err)
			return
		}

		assert.Equal(t, string(str), string(result))
		//t.Logf("str:%v buff:%v result:%v", string(str), string(buff), string(result))
	}

}
