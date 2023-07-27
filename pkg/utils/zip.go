package utils

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

func Zip(data []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)
	w := zlib.NewWriter(buffer)
	if err := WriteAll(w, data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Unzip(data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)

	reader, err := zlib.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	buff, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return buff, nil
}
