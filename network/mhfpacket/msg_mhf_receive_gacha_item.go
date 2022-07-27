package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfReceiveGachaItem represents the MSG_MHF_RECEIVE_GACHA_ITEM
type MsgMhfReceiveGachaItem struct{
	AckHandle      uint32
	Unk0           uint16
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfReceiveGachaItem) Opcode() network.PacketID {
	return network.MSG_MHF_RECEIVE_GACHA_ITEM
}

// Parse parses the packet from binary
func (m *MsgMhfReceiveGachaItem) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint16()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfReceiveGachaItem) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
