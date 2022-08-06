package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfGetBbsUserStatus represents the MSG_MHF_GET_BBS_USER_STATUS
type MsgMhfGetBbsUserStatus struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfGetBbsUserStatus) Opcode() network.PacketID {
	return network.MSG_MHF_GET_BBS_USER_STATUS
}

// Parse parses the packet from binary
func (m *MsgMhfGetBbsUserStatus) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgMhfGetBbsUserStatus) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
