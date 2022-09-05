package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfUpdateForceGuildRank represents the MSG_MHF_UPDATE_FORCE_GUILD_RANK
type MsgMhfUpdateForceGuildRank struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfUpdateForceGuildRank) Opcode() network.PacketID {
	return network.MSG_MHF_UPDATE_FORCE_GUILD_RANK
}

// Parse parses the packet from binary
func (m *MsgMhfUpdateForceGuildRank) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgMhfUpdateForceGuildRank) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
