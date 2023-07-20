package control

import "time"

type PeerInfo struct {
	PeerID         string
	ServerAddr     string
	StartTime      time.Time
	RegisterTime   time.Time
	LastAccessTime time.Time
}

type peerCache struct {
	curr      *PeerInfo
	peerInfos map[string]*PeerInfo
	connInfos map[string][]*PeerInfo
}

func newPeerCache() *peerCache {
	cache := new(peerCache)
	cache.curr = new(PeerInfo)
	cache.curr.StartTime = time.Now()
	cache.peerInfos = make(map[string]*PeerInfo)
	cache.connInfos = make(map[string][]*PeerInfo)
	return cache
}

func (cache *peerCache) SetPeerInfo(info *PeerInfo) {
	cache.peerInfos[info.PeerID] = info
}

func (cache *peerCache) SetConnPeerInfos(peerID string, infos []*PeerInfo) {
	cache.connInfos[peerID] = infos
}
