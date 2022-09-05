package mhfpacket

import (
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
	"erupe-ce/common/byteframe"
)

// MsgSysStageDestruct represents the MSG_SYS_STAGE_DESTRUCT
type MsgSysStageDestruct struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysStageDestruct) Opcode() network.PacketID {
	return network.MSG_SYS_STAGE_DESTRUCT
}

// Parse parses the packet from binary
func (m *MsgSysStageDestruct) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	// No data other than opcode
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgSysStageDestruct) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	// No data other than opcode
	return nil
}
