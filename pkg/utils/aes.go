package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// key.len 16 24 32
func AESEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	key = key[:block.BlockSize()]
	data = pkcs7Padding(data, len(key))
	crypted := make([]byte, len(data))
	cipher.NewCBCEncrypter(block, key).CryptBlocks(crypted, data)
	return crypted, nil
}

func AESDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	key = key[:block.BlockSize()]
	data := make([]byte, len(crypted))
	cipher.NewCBCDecrypter(block, key).CryptBlocks(data, crypted)
	return pkcs7UnPadding(data), nil
}

func pkcs7Padding(data []byte, size int) []byte {
	padding := size - len(data)%size
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
