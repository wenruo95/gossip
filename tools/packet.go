package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/wenruo95/gossip/pkg/log"
	"github.com/wenruo95/gossip/pkg/utils"
)

const (
	HeadFlag byte = 0x0a
	BodyFlag byte = 0x0b
	TailFlag byte = 0x0c

	CompressFlagGzip byte = 0x01
	EncodeFlagAES32  byte = 0x02

	DataSliceSize = 10 * 1024 * 1024
)

func EncodeFile(reader io.Reader, writer io.Writer, key []byte) error {
	keyMD5, err := utils.MD5Data(key)
	if err != nil {
		return err
	}

	txid := uint32(time.Now().Unix())
	head := &HeadPacket{CompressFlag: CompressFlagGzip, EncodeFlag: EncodeFlagAES32, EncKeySign: keyMD5}
	headBuff, err := head.Marshal()
	if err != nil {
		return err
	}
	if err := utils.WriteAll(writer, utils.Pack(headBuff, HeadFlag, txid)); err != nil {
		return err
	}

	var preMD5 []byte
	if m, err := utils.MD5Data(headBuff); err != nil {
		return err
	} else {
		preMD5 = m
		log.Debugf("[TEST1] md5:%v", hex.EncodeToString(preMD5))
	}

	size := DataSliceSize
	inputMD5 := md5.New()
	var sum int
	var seq uint32
	for {
		buff := make([]byte, size)
		n, err := reader.Read(buff)
		if err != nil && err != io.EOF {
			return err
		}
		sum = sum + n
		if _, err := inputMD5.Write(buff[:n]); err != nil {
			return err
		}

		buffMD5, err := utils.MD5Data(buff[:n])
		if err != nil {
			return err
		}

		encodeBuff, err := EncodeData(buff[:n], head.CompressFlag, head.EncodeFlag, key)
		if err != nil {
			return err
		}

		seq = seq + 1
		body := &BodyPacket{
			Seq:    seq,
			PreMD5: preMD5,
			MD5:    buffMD5,
			Data:   encodeBuff,
		}
		log.Debugf("[TEST] seq:%v pre:%v md5:%v data:%v", body.Seq, hex.EncodeToString(body.PreMD5), hex.EncodeToString(body.MD5), hex.EncodeToString(body.Data))
		bodyBuff, err := body.Marshal()
		if err != nil {
			return err
		}
		if m, err := utils.MD5Data(bodyBuff); err != nil {
			return err
		} else {
			preMD5 = m
			log.Debugf("[TEST2] md5:%v", hex.EncodeToString(preMD5))
		}
		if err := utils.WriteAll(writer, utils.Pack(bodyBuff, BodyFlag, txid)); err != nil {
			return err
		}
		if n < size {
			break
		}
	}

	imd5 := inputMD5.Sum(nil)
	tail := &TailPacket{PreMD5: preMD5, MD5: imd5}
	tailBuff, err := tail.Marshal()
	if err != nil {
		return err
	}
	if err := utils.WriteAll(writer, utils.Pack(tailBuff, TailFlag, txid)); err != nil {
		return err
	}

	log.Infof("length:%v md5:%v", sum, hex.EncodeToString(imd5))
	return nil
}

func DecodeFile(reader io.Reader, writer io.Writer, key []byte) error {
	keyMD5, err := utils.MD5Data(key)
	if err != nil {
		return err
	}

	var (
		head     *HeadPacket
		tail     *TailPacket
		sum      int
		preMD5   []byte
		seq      uint32
		checkMD5 = md5.New()
	)
	for idx := uint32(1); tail == nil; idx = idx + 1 {
		data, flag, txid, err := utils.Unpack(reader)
		log.Infof("flag:%v txid:%v sum:%v len:%v error:%v", flag, txid, sum, len(data), err)
		if err != nil && err != io.EOF {
			return err
		}

		if idx == 1 && flag != HeadFlag {
			return fmt.Errorf("invalid flag:%v expected:%v", flag, HeadFlag)
		}

		switch flag {
		case HeadFlag:
			if idx != 1 {
				return fmt.Errorf("invalid headflag idx:%v", idx)
			}
			head = new(HeadPacket)
			if err := head.Unmarshal(data); err != nil {
				return err
			}
			if !bytes.Equal(keyMD5, head.EncKeySign) {
				return fmt.Errorf("sign:%v error. expected:%v",
					hex.EncodeToString(keyMD5), hex.EncodeToString(head.EncKeySign))
			}

		case BodyFlag:
			seq = seq + 1
			body := new(BodyPacket)
			if err := body.Unmarshal(data); err != nil {
				return err
			}
			if body.Seq != seq {
				return fmt.Errorf("invalid seq:%v expected:%v", body.Seq, seq)
			}
			log.Debugf("[TEST] seq:%v pre:%v md5:%v data:%v", body.Seq, hex.EncodeToString(body.PreMD5), hex.EncodeToString(body.MD5), hex.EncodeToString(body.Data))
			if !bytes.Equal(preMD5, body.PreMD5) {
				return fmt.Errorf("pre sign:%v error. expected:%v",
					hex.EncodeToString(keyMD5), hex.EncodeToString(body.PreMD5))
			}
			buff, err := DecodeData(body.Data, head.CompressFlag, head.EncodeFlag, key)
			if err != nil {
				return err
			}
			m, err := utils.MD5Data(buff)
			if err != nil {
				return err
			}
			if !bytes.Equal(m, body.MD5) {
				return fmt.Errorf("current_data sign:%v error. expected:%v",
					hex.EncodeToString(m), hex.EncodeToString(body.MD5))
			}
			if _, err := checkMD5.Write(buff); err != nil {
				return err
			}
			sum = sum + len(buff)
			if err := utils.WriteAll(writer, buff); err != nil {
				return err
			}

		case TailFlag:
			tail = new(TailPacket)
			if err := tail.Unmarshal(data); err != nil {
				return err
			}
			if !bytes.Equal(preMD5, tail.PreMD5) {
				return fmt.Errorf("sign:%v error. expected:%v",
					hex.EncodeToString(keyMD5), hex.EncodeToString(tail.PreMD5))
			}
			break

		default:
		}

		m, err := utils.MD5Data(data)
		if err != nil {
			return err
		}
		preMD5 = m
		log.Debugf("[TEST1] md5:%v", hex.EncodeToString(preMD5))

	}

	fileMD5 := checkMD5.Sum(nil)
	if !bytes.Equal(fileMD5, tail.MD5) {
		return fmt.Errorf("file.md5:%v error, expected:%v",
			hex.EncodeToString(fileMD5), hex.EncodeToString(tail.MD5))
	}
	log.Infof("length:%v md5:%v", sum, hex.EncodeToString(fileMD5))
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
	if err := utils.WriteAll(buff, h.EncKeySign); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (h *HeadPacket) Unmarshal(body []byte) error {
	if len(body) != 18 {
		return fmt.Errorf("invalid body.len:%v expected:%v", len(body), 18)
	}
	h.CompressFlag = body[0]
	h.EncodeFlag = body[1]
	h.EncKeySign = body[2:]
	return nil
}

type BodyPacket struct {
	Seq    uint32 // 32
	PreMD5 []byte // 16
	MD5    []byte // 16
	Data   []byte // len - 64
}

func (b *BodyPacket) Marshal() ([]byte, error) {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, b.Seq)
	if err := utils.WriteAll(buff, b.PreMD5); err != nil {
		return nil, err
	}
	if err := utils.WriteAll(buff, b.MD5); err != nil {
		return nil, err
	}
	if err := utils.WriteAll(buff, b.Data); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (b *BodyPacket) Unmarshal(body []byte) error {
	if len(body) < 64 {
		return fmt.Errorf("invalid body.len:%v greater than %v", len(body), 64)
	}

	binary.Read(bytes.NewReader(body[:32]), binary.BigEndian, &b.Seq)
	b.PreMD5 = body[32:48]
	b.MD5 = body[48:64]
	b.Data = body[64:]
	return nil
}

type TailPacket struct {
	PreMD5 []byte // 16
	MD5    []byte // 16
}

func (t *TailPacket) Marshal() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := utils.WriteAll(buff, t.PreMD5); err != nil {
		return nil, err
	}
	if err := utils.WriteAll(buff, t.MD5); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
func (t *TailPacket) Unmarshal(body []byte) error {
	if len(body) < 12 {
		return fmt.Errorf("invalid body.len:%v greater than %v", len(body), 12)
	}
	t.PreMD5 = body[:16]
	t.MD5 = body[16:32]
	return nil
}

func EncodeData(data []byte, compressFlag, encodeFlag byte, key []byte) ([]byte, error) {
	if compressFlag == CompressFlagGzip {
		buff, err := utils.Zip(data)
		if err != nil {
			return nil, err
		}
		data = buff
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

func DecodeData(data []byte, compressFlag, encodeFlag byte, key []byte) ([]byte, error) {
	if encodeFlag == EncodeFlagAES32 {
		buff, err := utils.AESDecrypt(data, key)
		if err != nil {
			return nil, err
		}
		data = buff
	}

	if compressFlag == CompressFlagGzip {
		buff, err := utils.Unzip(data)
		if err != nil {
			return nil, err
		}
		data = buff
	}

	return data, nil
}
