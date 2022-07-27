package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysReserve71 represents the MSG_SYS_reserve71
type MsgSysReserve71 struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReserve71) Opcode() network.PacketID {
	return network.MSG_SYS_reserve71
}

// Parse parses the packet from binary
func (m *MsgSysReserve71) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysReserve71) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
