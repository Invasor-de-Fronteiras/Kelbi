package channelserver

import (
	"encoding/base64"
	"fmt"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
)

func handleMsgSysInsertUser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDeleteUser(s *Session, p mhfpacket.MHFPacket) {}

func broadcastNewUser(s *Session) {
	s.logger.Debug(fmt.Sprintf("Broadcasting new user: %s (%d)", s.Name, s.charID))

	clientNotif := byteframe.NewByteFrame()
	var temp mhfpacket.MHFPacket
	for _, session := range s.server.sessions {
		if session == s || !session.binariesDone {
			continue
		}
		temp = &mhfpacket.MsgSysInsertUser{CharID: session.charID}
		clientNotif.WriteUint16(uint16(temp.Opcode()))
		temp.Build(clientNotif, s.clientContext)
		for i := 0; i < 3; i++ {
			temp = &mhfpacket.MsgSysNotifyUserBinary{
				CharID:     session.charID,
				BinaryType: uint8(i + 1),
			}
			clientNotif.WriteUint16(uint16(temp.Opcode()))
			temp.Build(clientNotif, s.clientContext)
		}
	}
	s.QueueSend(clientNotif.Data())

	serverNotif := byteframe.NewByteFrame()
	temp = &mhfpacket.MsgSysInsertUser{CharID: s.charID}
	serverNotif.WriteUint16(uint16(temp.Opcode()))
	temp.Build(serverNotif, s.clientContext)
	for i := 0; i < 3; i++ {
		temp = &mhfpacket.MsgSysNotifyUserBinary{
			CharID:     s.charID,
			BinaryType: uint8(i + 1),
		}
		serverNotif.WriteUint16(uint16(temp.Opcode()))
		temp.Build(serverNotif, s.clientContext)
	}
	for _, session := range s.server.sessions {
		if session == s || !session.binariesDone {
			continue
		}
		session.QueueSend(serverNotif.Data())
	}
}

func handleMsgSysSetUserBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetUserBinary)
	s.server.userBinaryPartsLock.Lock()
	s.server.userBinaryParts[userBinaryPartID{charID: s.charID, index: pkt.BinaryType}] = pkt.RawDataPayload
	s.server.userBinaryPartsLock.Unlock()

	// Insert user once all binary parts exist
	if !s.binariesDone {
		for i := 0; i < 3; i++ {
			_, exists := s.server.userBinaryParts[userBinaryPartID{charID: s.charID, index: uint8(i + 1)}]
			if !exists {
				return
			}
		}
		s.binariesDone = true
		broadcastNewUser(s)
		return
	}

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

	// If we can't get the real data, use a placeholder.
	if !ok {
		if pkt.BinaryType == 1 {
			// Stub name response with character ID
			resp.WriteBytes([]byte(fmt.Sprintf("CID%d", s.charID)))
			resp.WriteUint8(0) // NULL terminator.
		} else if pkt.BinaryType == 2 {
			data, err := base64.StdEncoding.DecodeString("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBn8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAwAAAAAAAAAAAAAABAAAAAAAAAAAAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==")
			if err != nil {
				panic(err)
			}
			resp.WriteBytes(data)
		} else if pkt.BinaryType == 3 {
			data, err := base64.StdEncoding.DecodeString("AQAAA2ea5P8ATgEA/wEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBn8AAAAAAAAAAAABAKAMAAAAAAAAAAAAACgAAAAAAAAAAAABAsQOAAAAAAAAAAABA6UMAAAAAAAAAAABBKAMAAAAAAAAAAABBToNAAAAAAAAAAAAAAMAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAgACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/////AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
			if err != nil {
				panic(err)
			}
			resp.WriteBytes(data)
		}
	} else {
		resp.WriteBytes(data)
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgSysNotifyUserBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetBbsUserStatus(s *Session, p mhfpacket.MHFPacket) {}
