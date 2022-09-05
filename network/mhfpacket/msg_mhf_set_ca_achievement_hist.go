package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfSetCaAchievementHist represents the MSG_MHF_SET_CA_ACHIEVEMENT_HIST
type MsgMhfSetCaAchievementHist struct{
	AckHandle uint32
	Unk0 uint32
	Unk1 uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfSetCaAchievementHist) Opcode() network.PacketID {
	return network.MSG_MHF_SET_CA_ACHIEVEMENT_HIST
}

// Parse parses the packet from binary
func (m *MsgMhfSetCaAchievementHist) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint32()
	m.Unk1 = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfSetCaAchievementHist) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
