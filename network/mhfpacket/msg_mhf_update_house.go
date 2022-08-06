package mhfpacket

import (
	"errors"
	"erupe-ce/common/stringsupport"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfUpdateHouse represents the MSG_MHF_UPDATE_HOUSE
type MsgMhfUpdateHouse struct {
	AckHandle uint32
	State     uint8
	Unk1      uint8  // Always 0x01
	Unk2      uint16 // Always 0x0000
	Password  string
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfUpdateHouse) Opcode() network.PacketID {
	return network.MSG_MHF_UPDATE_HOUSE
}

// Parse parses the packet from binary
func (m *MsgMhfUpdateHouse) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.State = bf.ReadUint8()
	m.Unk1 = bf.ReadUint8()
	m.Unk2 = bf.ReadUint16()
	_ = bf.ReadUint8() // Password length
	m.Password = stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes())
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfUpdateHouse) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
