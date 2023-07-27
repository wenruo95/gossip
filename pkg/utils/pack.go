package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	DefaultPkgBegin   = 0x12
	DefaultPkgEnd     = 0x34
	DefaultPkgVersion = 0x01
)

var stdpacker = NewPacker(DefaultPkgBegin, DefaultPkgEnd, DefaultPkgVersion)
var Pack = stdpacker.Pack
var Unpack = stdpacker.Unpack

type Packer interface {
	Pack(body []byte, messageFlag byte, txid uint32) []byte
	Unpack(reader io.Reader) ([]byte, byte, uint32, error)
}

type packer struct {
	version  byte
	pkgBegin byte
	pkgEnd   byte
}

func NewPacker(pkgBegin, pkgEnd, version byte) *packer {
	p := new(packer)
	p.pkgBegin = pkgBegin
	p.pkgEnd = pkgEnd
	p.version = version
	return p
}

// 封包
// bodylen=len(body)
// 0--------1--------2--------3--------7--------11--------X--------X+1
// | start  | version| flags  | txid   | bodylen | body   | end    |
// +--------+--------+--------+--------+---------+--------+--------+
func (p *packer) Pack(body []byte, messageFlag byte, txid uint32) []byte {
	bodyLen := uint32(len(body))

	buf := &bytes.Buffer{}
	buf.WriteByte(p.pkgBegin)
	buf.WriteByte(p.version)
	buf.WriteByte(messageFlag)
	binary.Write(buf, binary.BigEndian, txid)
	binary.Write(buf, binary.BigEndian, bodyLen)
	buf.Write(body)
	buf.WriteByte(p.pkgEnd)

	return buf.Bytes()
}

// 解包 body messageflag txid err
func (p *packer) Unpack(reader io.Reader) ([]byte, byte, uint32, error) {
	const PkgHeadLen = 11
	headerData := make([]byte, PkgHeadLen)
	if _, err := io.ReadFull(reader, headerData); err != nil {
		return nil, 0, 0, err
	}

	headerBuf := bytes.NewBuffer(headerData)
	begin, _ := headerBuf.ReadByte()
	if begin != p.pkgBegin {
		return nil, 0, 0, fmt.Errorf("invalid pkghead:%v must be:%v", begin, p.pkgBegin)
	}

	version, err := headerBuf.ReadByte()
	if err != nil {
		return nil, 0, 0, err
	}
	if version != p.version {
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

	if data[len(data)-1] != p.pkgEnd {
		return nil, 0, 0, fmt.Errorf("invalid pkgend:%v must be:%v", data[len(data)-1], p.pkgEnd)
	}

	return data[:len(data)-1], messageFlag, txid, nil
}
