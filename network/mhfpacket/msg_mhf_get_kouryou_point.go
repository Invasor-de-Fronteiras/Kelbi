package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfGetKouryouPoint represents the MSG_MHF_GET_KOURYOU_POINT
type MsgMhfGetKouryouPoint struct {
	AckHandle uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfGetKouryouPoint) Opcode() network.PacketID {
	return network.MSG_MHF_GET_KOURYOU_POINT
}

// Parse parses the packet from binary
func (m *MsgMhfGetKouryouPoint) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfGetKouryouPoint) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
