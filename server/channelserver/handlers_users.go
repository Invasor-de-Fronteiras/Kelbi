package channelserver

import (
	"fmt"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
)

func handleMsgSysInsertUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetUserBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetUserBinary)
	s.server.userBinaryPartsLock.Lock()
	s.server.userBinaryParts[userBinaryPartID{charID: s.charID, index: pkt.BinaryType}] = pkt.RawDataPayload
	s.server.userBinaryPartsLock.Unlock()

	err := s.server.db.QueryRow("SELECT type1 FROM user_binaries WHERE id=$1", s.charID)
	if err != nil {
		s.server.db.Exec("INSERT INTO user_binaries (id) VALUES ($1)", s.charID)
	}

	s.server.db.Exec(fmt.Sprintf("UPDATE user_binaries SET type%d=$1 WHERE id=$2", pkt.BinaryType), pkt.RawDataPayload, s.charID)

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
	resp := byteframe.NewByteFrame()

	// If we can't get the real data, try to get it from the database.
	if !ok {
		var data []byte
		rows, _ := s.server.db.Queryx(fmt.Sprintf("SELECT type%d FROM user_binaries WHERE id=$1", pkt.BinaryType), pkt.CharID)
		for rows.Next() {
			rows.Scan(&data)
			resp.WriteBytes(data)
			doAckBufSucceed(s, pkt.AckHandle, resp.Data())
			return
		}
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		return
	} else {
		resp.WriteBytes(data)
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgSysNotifyUserBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBbsUserStatus(s *Session, p mhfpacket.MHFPacket) {}
