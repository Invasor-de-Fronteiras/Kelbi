package mhfpacket

import (
 "errors"

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfGetUdTacticsRanking represents the MSG_MHF_GET_UD_TACTICS_RANKING
type MsgMhfGetUdTacticsRanking struct {
  AckHandle uint32
  GuildID uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfGetUdTacticsRanking) Opcode() network.PacketID {
	return network.MSG_MHF_GET_UD_TACTICS_RANKING
}

// Parse parses the packet from binary
func (m *MsgMhfGetUdTacticsRanking) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
  m.AckHandle = bf.ReadUint32()
  m.GuildID = bf.ReadUint32()
  return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfGetUdTacticsRanking) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
