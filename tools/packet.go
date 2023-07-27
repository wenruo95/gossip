package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io"

	"github.com/wenruo95/gossip/pkg/log"
	"github.com/wenruo95/gossip/pkg/tcp"
	"github.com/wenruo95/gossip/pkg/utils"
)

const (
	HeadFlag byte = 0x01
	BodyFlag byte = 0x02
	TailFlag byte = 0x03

	CompressFlagGzip byte = 0x01

	EncodeFlagAES32 byte = 0x02

	DataSliceSize = 10 * 1024 * 1024
)

func EncodeFile(reader io.Reader, writer io.Writer, key []byte) error {
	keyMD5, err := utils.MD5Data(key)
	if err != nil {
		return err
	}

	var seq uint32 = 1
	head := &HeadPacket{CompressFlag: CompressFlagGzip, EncodeFlag: EncodeFlagAES32, EncKeySign: keyMD5}
	headBuff, err := head.Marshal()
	if err != nil {
		return err
	}
	if err := tcp.WriteAll(writer, tcp.Pack(headBuff, HeadFlag, seq)); err != nil {
		return err
	}

	var preMD5 []byte
	if m, err := utils.MD5Data(headBuff); err != nil {
		return err
	} else {
		preMD5 = m
	}

	size := DataSliceSize
	inputMD5 := md5.New()
	var sum int
	for {
		buff := make([]byte, size)
		n, err := reader.Read(buff)
		if err != nil && err != io.EOF {
			return err
		}
		if _, err := inputMD5.Write(buff[:n]); err != nil {
			return err
		}

		localMD5, err := utils.MD5Data(buff[:n])
		if err != nil {
			return err
		}

		encodeBuff, err := EncodeData(buff[:n], head.CompressFlag, head.EncodeFlag, key)
		if err != nil {
			return err
		}

		body := &BodyPacket{
			PreMD5: preMD5,
			MD5:    localMD5,
			Data:   encodeBuff,
		}
		data, err := body.Marshal()
		if err != nil {
			return err
		}
		seq = seq + 1
		if err := tcp.WriteAll(writer, tcp.Pack(data, BodyFlag, seq)); err != nil {
			return err
		}
		sum = sum + n
		preMD5 = localMD5

		if n < size {
			break
		}
	}

	seq = seq + 1
	imd5 := inputMD5.Sum(nil)
	tail := &TailPacket{PreMD5: preMD5, MD5: imd5}
	tailBuff, err := tail.Marshal()
	if err != nil {
		return err
	}
	if err := tcp.WriteAll(writer, tcp.Pack(tailBuff, TailFlag, seq)); err != nil {
		return err
	}

	log.Infof("length:%v md5:%v", sum, hex.EncodeToString(imd5))
	return nil
}

func DecodeFile(reader io.Reader, writer io.Writer, key []byte) error {
	return nil
}

type HeadPacket struct {
	CompressFlag byte   // gzip
	EncodeFlag   byte   // aes
	EncKeySign   []byte // encode_key_pre
}

func (h *HeadPacket) Marshal() ([]byte, error) {
	buff := new(bytes.Buffer)
	buff.WriteByte(h.CompressFlag)
	buff.WriteByte(h.EncodeFlag)
	if err := tcp.WriteAll(buff, h.EncKeySign); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

type BodyPacket struct {
	PreMD5 []byte // 16
	MD5    []byte // 16
	Data   []byte // len - 32
}

func (b *BodyPacket) Marshal() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := tcp.WriteAll(buff, b.PreMD5); err != nil {
		return nil, err
	}
	if err := tcp.WriteAll(buff, b.MD5); err != nil {
		return nil, err
	}
	if err := tcp.WriteAll(buff, b.Data); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

type TailPacket struct {
	PreMD5 []byte // 16
	MD5    []byte // 16
}

func (t *TailPacket) Marshal() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := tcp.WriteAll(buff, t.PreMD5); err != nil {
		return nil, err
	}
	if err := tcp.WriteAll(buff, t.MD5); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func EncodeData(data []byte, compressFlag, encodeFlag byte, key []byte) ([]byte, error) {
	if compressFlag == CompressFlagGzip {
		buffer := new(bytes.Buffer)
		w := gzip.NewWriter(buffer)
		if _, err := w.Write(data); err != nil {
			return nil, err
		}
		w.Close()
		data = buffer.Bytes()
	}
	if encodeFlag == EncodeFlagAES32 {
		buff, err := utils.AESEncrypt(data, key)
		if err != nil {
			return nil, err
		}
		data = buff
	}

	return data, nil
}
