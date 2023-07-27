package tcp

import "github.com/wenruo95/gossip/pkg/utils"

const (
	PkgBegin   = 0xEF
	PkgEnd     = 0xFF
	PkgVersion = 0x22
	PkgHeadLen = 11
)

var stdpacker = utils.NewPacker(PkgBegin, PkgEnd, PkgVersion)
var Pack = stdpacker.Pack
var Unpack = stdpacker.Unpack
var WriteAll = utils.WriteAll
