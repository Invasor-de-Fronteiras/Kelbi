package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysReserve0E represents the MSG_SYS_reserve0E
type MsgSysReserve0E struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReserve0E) Opcode() network.PacketID {
	return network.MSG_SYS_reserve0E
}

// Parse parses the packet from binary
func (m *MsgSysReserve0E) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysReserve0E) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
