package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysReserve4A represents the MSG_SYS_reserve4A
type MsgSysReserve4A struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReserve4A) Opcode() network.PacketID {
	return network.MSG_SYS_reserve4A
}

// Parse parses the packet from binary
func (m *MsgSysReserve4A) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysReserve4A) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
