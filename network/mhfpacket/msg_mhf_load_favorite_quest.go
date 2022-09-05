package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfLoadFavoriteQuest represents the MSG_MHF_LOAD_FAVORITE_QUEST
type MsgMhfLoadFavoriteQuest struct {
	AckHandle uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfLoadFavoriteQuest) Opcode() network.PacketID {
	return network.MSG_MHF_LOAD_FAVORITE_QUEST
}

// Parse parses the packet from binary
func (m *MsgMhfLoadFavoriteQuest) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfLoadFavoriteQuest) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
