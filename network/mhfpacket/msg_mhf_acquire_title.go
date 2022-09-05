package mhfpacket

import (
	"errors"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfAcquireTitle represents the MSG_MHF_ACQUIRE_TITLE
type MsgMhfAcquireTitle struct {
	AckHandle uint32
	Unk0      uint16
	Unk1      uint16
	TitleID   uint16
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfAcquireTitle) Opcode() network.PacketID {
	return network.MSG_MHF_ACQUIRE_TITLE
}

// Parse parses the packet from binary
func (m *MsgMhfAcquireTitle) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk0 = bf.ReadUint16()
	m.Unk1 = bf.ReadUint16()
	m.TitleID = bf.ReadUint16()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfAcquireTitle) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
