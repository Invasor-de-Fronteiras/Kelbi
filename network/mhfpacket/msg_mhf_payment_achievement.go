package mhfpacket

import (
	"errors"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfPaymentAchievement represents the MSG_MHF_PAYMENT_ACHIEVEMENT
type MsgMhfPaymentAchievement struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfPaymentAchievement) Opcode() network.PacketID {
	return network.MSG_MHF_PAYMENT_ACHIEVEMENT
}

// Parse parses the packet from binary
func (m *MsgMhfPaymentAchievement) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgMhfPaymentAchievement) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
