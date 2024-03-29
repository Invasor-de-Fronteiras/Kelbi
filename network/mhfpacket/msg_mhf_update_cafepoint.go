package mhfpacket

import (
	"errors"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfUpdateCafepoint represents the MSG_MHF_UPDATE_CAFEPOINT
type MsgMhfUpdateCafepoint struct {
	AckHandle uint32
	Unk0      uint16 // Hardcoded 0 in binary
	Unk1      uint16 // Hardcoded 0 in binary
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfUpdateCafepoint) Opcode() network.PacketID {
	return network.MSG_MHF_UPDATE_CAFEPOINT
}

// Parse parses the packet from binary
func (m *MsgMhfUpdateCafepoint) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint16()
	m.Unk1 = bf.ReadUint16()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfUpdateCafepoint) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
