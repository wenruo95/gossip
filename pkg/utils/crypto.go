package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func MD5Data(data []byte) ([]byte, error) {
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func MD5File(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func Sha1Data(data []byte) ([]byte, error) {
	h := sha1.New()
	n, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, fmt.Errorf("invalid n:%v expected:%v", n, len(data))
	}
	return h.Sum(nil), nil
}
