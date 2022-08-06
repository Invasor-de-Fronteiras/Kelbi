package mhfpacket

import (
 "errors"

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysReserve1A4 represents the MSG_SYS_reserve1A4
type MsgSysReserve1A4 struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReserve1A4) Opcode() network.PacketID {
	return network.MSG_SYS_reserve1A4
}

// Parse parses the packet from binary
func (m *MsgSysReserve1A4) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysReserve1A4) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
