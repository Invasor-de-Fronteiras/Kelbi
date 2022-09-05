package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfGetExtraInfo represents the MSG_MHF_GET_EXTRA_INFO
type MsgMhfGetExtraInfo struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfGetExtraInfo) Opcode() network.PacketID {
	return network.MSG_MHF_GET_EXTRA_INFO
}

// Parse parses the packet from binary
func (m *MsgMhfGetExtraInfo) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgMhfGetExtraInfo) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
