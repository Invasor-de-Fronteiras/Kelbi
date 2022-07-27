package mhfpacket

import (
	"errors"
	"erupe-ce/common/stringsupport"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfEnumerateHouse represents the MSG_MHF_ENUMERATE_HOUSE
type MsgMhfEnumerateHouse struct {
	AckHandle uint32
	CharID    uint32
	Method    uint8
	Unk       uint16
	Name      string
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfEnumerateHouse) Opcode() network.PacketID {
	return network.MSG_MHF_ENUMERATE_HOUSE
}

// Parse parses the packet from binary
func (m *MsgMhfEnumerateHouse) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.CharID = bf.ReadUint32()
	m.Method = bf.ReadUint8()
	m.Unk = bf.ReadUint16()
	_ = bf.ReadUint8() // len
	m.Name = stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes())
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfEnumerateHouse) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
