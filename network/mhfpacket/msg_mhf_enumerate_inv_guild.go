package mhfpacket

import (
	"errors"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfEnumerateInvGuild represents the MSG_MHF_ENUMERATE_INV_GUILD
type MsgMhfEnumerateInvGuild struct {
	AckHandle uint32
	Unk       []byte
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfEnumerateInvGuild) Opcode() network.PacketID {
	return network.MSG_MHF_ENUMERATE_INV_GUILD
}

// Parse parses the packet from binary
func (m *MsgMhfEnumerateInvGuild) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.Unk = bf.ReadBytes(9)
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfEnumerateInvGuild) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
