package channelserver

import (
	"erupe-ce/network/mhfpacket"
)

func handleMsgSysInsertUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetUserBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetUserBinary)
	s.server.userBinaryPartsLock.Lock()
	s.server.userBinaryParts[userBinaryPartID{charID: s.charID, index: pkt.BinaryType}] = pkt.RawDataPayload
	s.server.userBinaryPartsLock.Unlock()

	msg := &mhfpacket.MsgSysNotifyUserBinary{
		CharID:     s.charID,
		BinaryType: pkt.BinaryType,
	}

	s.server.BroadcastMHF(msg, s)
}

func handleMsgSysGetUserBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysGetUserBinary)

	// Try to get the data.
	s.server.userBinaryPartsLock.RLock()
	defer s.server.userBinaryPartsLock.RUnlock()
	data, ok := s.server.userBinaryParts[userBinaryPartID{charID: pkt.CharID, index: pkt.BinaryType}]
	// resp := byteframe.NewByteFrame()

	// If we can't get the real data, fail.
	if !ok {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
	} else {
		doAckBufSucceed(s, pkt.AckHandle, data)
	}
}

func handleMsgSysNotifyUserBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBbsUserStatus(s *Session, p mhfpacket.MHFPacket) {}
