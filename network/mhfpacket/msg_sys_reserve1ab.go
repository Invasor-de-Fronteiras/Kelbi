package mhfpacket

import (
 "errors"

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysReserve1AB represents the MSG_SYS_reserve1AB
type MsgSysReserve1AB struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReserve1AB) Opcode() network.PacketID {
	return network.MSG_SYS_reserve1AB
}

// Parse parses the packet from binary
func (m *MsgSysReserve1AB) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysReserve1AB) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
