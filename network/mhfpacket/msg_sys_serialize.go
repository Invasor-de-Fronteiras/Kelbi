package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysSerialize represents the MSG_SYS_SERIALIZE
type MsgSysSerialize struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysSerialize) Opcode() network.PacketID {
	return network.MSG_SYS_SERIALIZE
}

// Parse parses the packet from binary
func (m *MsgSysSerialize) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysSerialize) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
