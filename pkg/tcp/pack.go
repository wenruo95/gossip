package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	PkgBegin   = 0xEF
	PkgEnd     = 0xFF
	PkgVersion = 0x01
	PkgHeadLen = 11
)

// 封包
// bodylen=len(body)
// 0--------1--------2--------3--------7--------11--------X--------X+1
// | start  | version| flags  | txid   | bodylen | body   | end    |
// +--------+--------+--------+--------+---------+--------+--------+
func Pack(body []byte, messageFlag byte, txid uint32) []byte {
	bodyLen := uint32(len(body))

	buf := &bytes.Buffer{}
	buf.WriteByte(PkgBegin)
	buf.WriteByte(PkgVersion)
	buf.WriteByte(messageFlag)
	binary.Write(buf, binary.BigEndian, txid)
	binary.Write(buf, binary.BigEndian, bodyLen)
	buf.Write(body)
	buf.WriteByte(PkgEnd)

	return buf.Bytes()
}

// 解包 body messageflag txid err
func Unpack(reader io.Reader) ([]byte, byte, uint32, error) {
	headerData := make([]byte, PkgHeadLen)
	if _, err := io.ReadFull(reader, headerData); err != nil {
		return nil, 0, 0, err
	}

	headerBuf := bytes.NewBuffer(headerData)
	begin, _ := headerBuf.ReadByte()
	if begin != PkgBegin {
		return nil, 0, 0, fmt.Errorf("invalid pkghead:%v must be:%v", begin, PkgBegin)
	}

	version, err := headerBuf.ReadByte()
	if err != nil {
		return nil, 0, 0, err
	}
	if version != PkgVersion {
		return nil, 0, 0, fmt.Errorf("invalid version:%v", version)
	}

	messageFlag, err := headerBuf.ReadByte()
	if err != nil {
		return nil, 0, 0, err
	}

	var txid, size uint32
	binary.Read(headerBuf, binary.BigEndian, &txid)
	binary.Read(headerBuf, binary.BigEndian, &size)

	data := make([]byte, size+1)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, 0, 0, err
	}

	if data[len(data)-1] != PkgEnd {
		return nil, 0, 0, fmt.Errorf("invalid pkgend:%v must be:%v", data[len(data)-1], PkgEnd)
	}

	return data[:len(data)-1], messageFlag, txid, nil
}

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
