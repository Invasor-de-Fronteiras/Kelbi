package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfGetBoxGachaInfo represents the MSG_MHF_GET_BOX_GACHA_INFO
type MsgMhfGetBoxGachaInfo struct{
	AckHandle uint32
	GachaHash uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfGetBoxGachaInfo) Opcode() network.PacketID {
	return network.MSG_MHF_GET_BOX_GACHA_INFO
}

// Parse parses the packet from binary
func (m *MsgMhfGetBoxGachaInfo) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.GachaHash = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfGetBoxGachaInfo) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
