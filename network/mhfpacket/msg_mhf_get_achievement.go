package mhfpacket

import (
	"errors"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfGetAchievement represents the MSG_MHF_GET_ACHIEVEMENT
type MsgMhfGetAchievement struct {
	AckHandle uint32
	CharID    uint32
	Unk1      uint32 // char?
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfGetAchievement) Opcode() network.PacketID {
	return network.MSG_MHF_GET_ACHIEVEMENT
}

// Parse parses the packet from binary
func (m *MsgMhfGetAchievement) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.CharID = bf.ReadUint32()
	m.Unk1 = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfGetAchievement) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
