package utils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {

	type testData struct {
		dataLen int
		keyLen  int
	}

	datas := []*testData{
		{dataLen: 10, keyLen: 16},
		{dataLen: 10, keyLen: 24},
		{dataLen: 10, keyLen: 32},
		{dataLen: 50, keyLen: 16},
		{dataLen: 60, keyLen: 24},
		{dataLen: 300, keyLen: 32},
		{dataLen: 600, keyLen: 16},
		{dataLen: 900, keyLen: 24},
		{dataLen: 1000, keyLen: 32},
	}

	rand.Seed(time.Now().UnixNano())
	for _, testData := range datas {
		for i := 0; i < 10; i++ {
			data := RandCharacter(testData.dataLen)
			key := RandCharacter(testData.keyLen)

			crypted, err := AESEncrypt(data, key)
			assert.Nil(t, err, testData.dataLen, testData.keyLen)

			decryptData, err := AESDecrypt(crypted, key)
			assert.Nil(t, err, testData.dataLen, testData.keyLen)

			assert.Equal(t, data, decryptData)
		}
	}

}
