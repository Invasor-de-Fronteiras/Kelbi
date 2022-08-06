package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfGetBreakSeibatuLevelReward represents the MSG_MHF_GET_BREAK_SEIBATU_LEVEL_REWARD
type MsgMhfGetBreakSeibatuLevelReward struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfGetBreakSeibatuLevelReward) Opcode() network.PacketID {
	return network.MSG_MHF_GET_BREAK_SEIBATU_LEVEL_REWARD
}

// Parse parses the packet from binary
func (m *MsgMhfGetBreakSeibatuLevelReward) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgMhfGetBreakSeibatuLevelReward) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
