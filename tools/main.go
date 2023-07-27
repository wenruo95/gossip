package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/wenruo95/gossip/pkg/log"
	"github.com/wenruo95/gossip/pkg/utils"
)

var (
	inputFile  string
	outputFile string
	password   string
	op         string
)

const DefaultEncKeyPre = "20230727_zw_"

func generateEncKey(str string) ([]byte, error) {
	key, err := utils.MD5Data([]byte(str))
	if err != nil {
		return nil, err
	}
	data := make([]byte, 0)
	for i := 0; i < len(key); i++ {
		data = append(data, key[i], key[len(key)-i-1])
	}
	return data, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.StringVar(&inputFile, "i", "", "-i=./hello.test")
	flag.StringVar(&outputFile, "o", "", "-o=./hello.test.encrypted")
	flag.StringVar(&password, "p", "", "-p=12345678")
	flag.StringVar(&op, "op", "", "-op=enc/dec")
	flag.Parse()
}

func main() {
	if len(inputFile) == 0 || len(op) == 0 {
		log.Fatalf("invalid filepath:%v op:%v", inputFile, op)
	}
	if len(outputFile) == 0 {
		outputFile = fmt.Sprintf("%s.%s.%v", op, inputFile, time.Now().Unix())
	}

	inputFd, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("open file:%v error:%v", inputFile, err)
	}
	defer inputFd.Close()

	outputFd, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("open file:%v error:%v", outputFile, err)
	}
	defer outputFd.Close()

	str := DefaultEncKeyPre + password
	key, err := generateEncKey(str)
	if err != nil {
		log.Fatalf("generate_enc_key error:%v", err)
	}

	begin := time.Now()
	switch op {
	case "enc":
		if err := EncodeFile(inputFd, outputFd, key); err != nil {
			log.Fatalf("encode_file error:%v", err)
		}

	case "dec":
		if err := DecodeFile(inputFd, outputFd, key); err != nil {
			log.Fatalf("decode_file error:%v", err)
		}

	default:

	}

	log.Infof("finished op:%v consume:%v", op, time.Since(begin))
	return
}
