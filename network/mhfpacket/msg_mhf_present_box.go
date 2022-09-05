package mhfpacket

import (
	"errors"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfPresentBox represents the MSG_MHF_PRESENT_BOX
type MsgMhfPresentBox struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfPresentBox) Opcode() network.PacketID {
	return network.MSG_MHF_PRESENT_BOX
}

// Parse parses the packet from binary
func (m *MsgMhfPresentBox) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgMhfPresentBox) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
