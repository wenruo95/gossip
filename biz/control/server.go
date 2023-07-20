package control

import (
	"strconv"

	"github.com/wenruo95/gossip/pkg/log"
	"github.com/wenruo95/gossip/pkg/tcp"
)

// 1. register to server_addr list
// 2. heartbeat
// 3. query peer_info
// 4. query peer_meta_info

func (ctrl *control) OnConnect(conn *tcp.ClientConn) {
	log.Infof("server: new client connection. addr:%v", conn.RemoteAddr().String())
}

func (ctrl *control) OnMessage(conn *tcp.ClientConn, body []byte, messageFlag byte, txid uint32) {
	log.Infof("server: recieve message. len:%v flag:%v txid:%v body:%s", len(body), messageFlag, txid, string(body))
	if err := conn.Send([]byte("server_hello_world_"+strconv.FormatUint(uint64(txid), 10)), messageFlag, txid); err != nil {
		conn.Close()
		log.Errorf("error:%v", err)
	}
}

func (ctrl *control) OnDisconnect(conn *tcp.ClientConn, reason string) {
	log.Infof("server: client disconnected. addr:%v reason:%v", conn.RemoteAddr().String(), reason)
}
