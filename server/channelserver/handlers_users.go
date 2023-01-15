package channelserver

import (
	"erupe-ce/network/mhfpacket"
	"fmt"
)

func handleMsgSysInsertUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetUserBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetUserBinary)
	s.Server.userBinaryPartsLock.Lock()
	s.Server.userBinaryParts[userBinaryPartID{charID: s.CharID, index: pkt.BinaryType}] = pkt.RawDataPayload
	s.Server.userBinaryPartsLock.Unlock()

	var exists []byte
	err := s.Server.db.QueryRow("SELECT type2 FROM user_binary WHERE id=$1", s.CharID).Scan(&exists)
	if err != nil {
		// nolint:errcheck
		s.Server.db.Exec("INSERT INTO user_binary (id) VALUES ($1)", s.CharID)
	}

	// nolint:errcheck
	s.Server.db.Exec(fmt.Sprintf("UPDATE user_binary SET type%d=$1 WHERE id=$2", pkt.BinaryType), pkt.RawDataPayload, s.CharID)

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
