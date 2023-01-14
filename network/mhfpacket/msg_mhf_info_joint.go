package mhfpacket

import (
	"errors"

	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

// MsgMhfInfoJoint represents the MSG_MHF_INFO_JOINT
type MsgMhfInfoJoint struct {
	AckHandle  uint32
	AllianceID uint32
	Unk        uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfInfoJoint) Opcode() network.PacketID {
	return network.MSG_MHF_INFO_JOINT
}

// Parse parses the packet from binary
func (m *MsgMhfInfoJoint) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	m.AllianceID = bf.ReadUint32()
	m.Unk = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfInfoJoint) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
