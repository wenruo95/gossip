package tcp

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type TestTCPServerHandler struct {
}

func (handler *TestTCPServerHandler) OnConnect(conn *ClientConn) {
	log.Printf("[I] server: new client connection. addr:%v", conn.RemoteAddr().String())
}

func (handler *TestTCPServerHandler) OnMessage(conn *ClientConn, body []byte, messageFlag byte, txid uint32) {
	log.Printf("[I] server: recieve message. len:%v flag:%v txid:%v body:%s", len(body), messageFlag, txid, string(body))
	conn.Send([]byte("server_hello_world_"+strconv.FormatUint(uint64(txid), 10)), messageFlag, txid)
}

func (handler *TestTCPServerHandler) OnDisconnect(conn *ClientConn, reason string) {
	log.Printf("[I] server: client disconnected. addr:%v reason:%v", conn.RemoteAddr().String(), reason)
}

func Test_ServerTimeout(t *testing.T) {
	svr := NewServer(
		WithAddr("127.0.0.1:8000"),
		WithHandler(&TestTCPServerHandler{}),
		WithTimeout(10*time.Second),
	)
	go svr.Serve()
	time.Sleep(time.Second)
	svr.Close()
}
