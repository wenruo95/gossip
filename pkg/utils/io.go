package utils

import (
	"fmt"
	"io"
)

func WriteAll(writer io.Writer, data []byte) error {
	var size int
	for {
		body := data[size:]
		n, err := writer.Write(body)
		if err != nil {
			return err
		}

		size = size + n
		if size >= len(data) {
			break
		}
	}
	if size != len(data) {
		return fmt.Errorf("write error. n:%v expected:%v", size, len(data))
	}

	return nil
}
