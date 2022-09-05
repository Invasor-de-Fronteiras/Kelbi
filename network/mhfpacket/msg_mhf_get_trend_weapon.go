package mhfpacket

import (
	"errors"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfGetTrendWeapon represents the MSG_MHF_GET_TREND_WEAPON
type MsgMhfGetTrendWeapon struct {
	AckHandle uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfGetTrendWeapon) Opcode() network.PacketID {
	return network.MSG_MHF_GET_TREND_WEAPON
}

// Parse parses the packet from binary
func (m *MsgMhfGetTrendWeapon) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfGetTrendWeapon) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
