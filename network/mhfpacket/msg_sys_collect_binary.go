package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysCollectBinary represents the MSG_SYS_COLLECT_BINARY
type MsgSysCollectBinary struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysCollectBinary) Opcode() network.PacketID {
	return network.MSG_SYS_COLLECT_BINARY
}

// Parse parses the packet from binary
func (m *MsgSysCollectBinary) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysCollectBinary) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
