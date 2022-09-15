package channelserver

import (
	"erupe-ce/network/mhfpacket"
)

func handleMsgSysInsertUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetUserBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetUserBinary)
	s.Server.userBinaryPartsLock.Lock()
	s.Server.userBinaryParts[userBinaryPartID{charID: s.CharID, index: pkt.BinaryType}] = pkt.RawDataPayload
	s.Server.userBinaryPartsLock.Unlock()

	msg := &mhfpacket.MsgSysNotifyUserBinary{
		CharID:     s.CharID,
		BinaryType: pkt.BinaryType,
	}

	s.Server.BroadcastMHF(msg, s)
}

func handleMsgSysGetUserBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysGetUserBinary)

	// Try to get the data.
	s.Server.userBinaryPartsLock.RLock()
	defer s.Server.userBinaryPartsLock.RUnlock()
	data, ok := s.Server.userBinaryParts[userBinaryPartID{charID: pkt.CharID, index: pkt.BinaryType}]
	// resp := byteframe.NewByteFrame()

	// If we can't get the real data, fail.
	if !ok {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
	} else {
		doAckBufSucceed(s, pkt.AckHandle, data)
	}
}

func handleMsgSysNotifyUserBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBbsUserStatus(s *Session, p mhfpacket.MHFPacket) {}
