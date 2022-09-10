package channelserver

import (
	"encoding/binary"
	"encoding/hex"
	"erupe-ce/common/stringsupport"
	"fmt"
	"io"
	"math"
	"net"
	"strings"
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
	"math/bits"
	"math/rand"

	"go.uber.org/zap"
)

// Temporary function to just return no results for a MSG_MHF_ENUMERATE* packet
func stubEnumerateNoResults(s *Session, ackHandle uint32) {
	enumBf := byteframe.NewByteFrame()
	enumBf.WriteUint32(0) // Entry count (count for quests, rankings, events, etc.)

	doAckBufSucceed(s, ackHandle, enumBf.Data())
}

// Temporary function to just return no results for many MSG_MHF_GET* packets.
func stubGetNoResults(s *Session, ackHandle uint32) {
	resp := byteframe.NewByteFrame()
	resp.WriteUint32(0x0A218EAD) // Unk shared ID. Sent in response of MSG_MHF_GET_TOWER_INFO, MSG_MHF_GET_PAPER_DATA etc. (World ID?)
	resp.WriteUint32(0)          // Unk
	resp.WriteUint32(0)          // Unk
	resp.WriteUint32(0)          // Entry count

	doAckBufSucceed(s, ackHandle, resp.Data())
}

func doAckBufSucceed(s *Session, ackHandle uint32, data []byte) {
	s.QueueSendMHF(&mhfpacket.MsgSysAck{
		AckHandle:        ackHandle,
		IsBufferResponse: true,
		ErrorCode:        0,
		AckData:          data,
	})
}

func doAckBufFail(s *Session, ackHandle uint32, data []byte) {
	s.QueueSendMHF(&mhfpacket.MsgSysAck{
		AckHandle:        ackHandle,
		IsBufferResponse: true,
		ErrorCode:        1,
		AckData:          data,
	})
}

func doAckSimpleSucceed(s *Session, ackHandle uint32, data []byte) {
	s.QueueSendMHF(&mhfpacket.MsgSysAck{
		AckHandle:        ackHandle,
		IsBufferResponse: false,
		ErrorCode:        0,
		AckData:          data,
	})
}

func doAckSimpleFail(s *Session, ackHandle uint32, data []byte) {
	s.QueueSendMHF(&mhfpacket.MsgSysAck{
		AckHandle:        ackHandle,
		IsBufferResponse: false,
		ErrorCode:        1,
		AckData:          data,
	})
}

func updateRights(s *Session) {
	s.Rights = uint32(0x0E)
	err := s.Server.db.QueryRow("SELECT rights FROM users u INNER JOIN characters c ON u.id = c.user_id WHERE c.id = $1", s.CharID).Scan(&s.Rights)
	if err != nil {
		s.logger.Fatal("FAILED UPDATE RIGHTS", zap.Error(err))
	}

	rights := make([]mhfpacket.ClientRight, 0)
	tempRights := s.Rights
	for i := 30; i > 0; i-- {
		right := uint32(math.Pow(2, float64(i)))
		if tempRights-right < 0x80000000 {
			if i == 1 {
				continue
			}
			rights = append(rights, mhfpacket.ClientRight{ID: uint16(i), Timestamp: 0x70DB59F0})
			tempRights -= right
		}
	}
	rights = append(rights, mhfpacket.ClientRight{ID: 1, Timestamp: 0})

	update := &mhfpacket.MsgSysUpdateRight{
		ClientRespAckHandle: 0,
		Bitfield:            s.Rights,
		Rights:              rights,
		UnkSize:             0,
	}
	s.QueueSendMHF(update)
}

func handleMsgHead(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysExtendThreshold(s *Session, p mhfpacket.MHFPacket) {
	// No data aside from header, no resp required.
}

func handleMsgSysEnd(s *Session, p mhfpacket.MHFPacket) {
	// No data aside from header, no resp required.
}

func handleMsgSysNop(s *Session, p mhfpacket.MHFPacket) {
	// No data aside from header, no resp required.
}

func handleMsgSysAck(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysTerminalLog(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysTerminalLog)
	resp := byteframe.NewByteFrame()

	resp.WriteUint32(0x98bd51a9) // LogID to use for requests after this.
	doAckSimpleSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgSysLogin(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysLogin)

	if !s.Server.erupeConfig.DevModeOptions.DisableTokenCheck {
		var token string
		err := s.Server.db.QueryRow("SELECT token FROM sign_sessions WHERE token=$1", pkt.LoginTokenString).Scan(&token)
		if err != nil {
			s.rawConn.Close()
			s.logger.Warn(fmt.Sprintf("Invalid login token, offending CID: (%d)", pkt.CharID0))
			return
		}
	}

	s.Lock()
	s.CharID = pkt.CharID0
	s.Token = pkt.LoginTokenString
	s.Unlock()
	updateRights(s)

	bf := byteframe.NewByteFrame()
	bf.WriteUint32(uint32(Time_Current_Adjusted().Unix())) // Unix timestamp

	_, err := s.Server.db.Exec("UPDATE servers SET current_players=$1 WHERE server_id=$2", len(s.Server.Sessions), s.Server.ID)
	if err != nil {
		panic(err)
	}

	_, err = s.Server.db.Exec("UPDATE sign_sessions SET server_id=$1, char_id=$2 WHERE token=$3", s.Server.ID, s.CharID, s.Token)
	if err != nil {
		panic(err)
	}

	_, err = s.Server.db.Exec("UPDATE sign_sessions SET server_id=$1, char_id=$2 WHERE token=$3", s.Server.ID, s.CharID, s.Token)
	if err != nil {
		panic(err)
	}

	_, err = s.Server.db.Exec("UPDATE characters SET last_login=$1 WHERE id=$2", Time_Current().Unix(), s.CharID)
	if err != nil {
		panic(err)
	}

	_, err = s.Server.db.Exec("UPDATE users u SET last_character=$1 WHERE u.id=(SELECT c.user_id FROM characters c WHERE c.id=$1)", s.CharID)
	if err != nil {
		panic(err)
	}

	_, err = s.Server.db.Exec("UPDATE users u SET last_character=$1 WHERE u.id=(SELECT c.user_id FROM characters c WHERE c.id=$1)", s.CharID)
	if err != nil {
		panic(err)
	}

	doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())

	updateRights(s)

	s.Server.BroadcastMHF(&mhfpacket.MsgSysInsertUser{CharID: s.CharID}, s)
}

func handleMsgSysLogout(s *Session, p mhfpacket.MHFPacket) {
	logoutPlayer(s)
}

func logoutPlayer(s *Session) {
	if _, exists := s.Server.Sessions[s.rawConn]; exists {
		delete(s.Server.Sessions, s.rawConn)
		s.rawConn.Close()
	} else {
		return // Prevent re-running logout logic on real logouts
	}

	_, err := s.Server.db.Exec("UPDATE sign_sessions SET server_id=NULL, char_id=NULL WHERE token=$1", s.Token)
	if err != nil {
		panic(err)
	}

	_, err = s.Server.db.Exec("UPDATE servers SET current_players=$1 WHERE server_id=$2", len(s.Server.Sessions), s.Server.ID)
	if err != nil {
		panic(err)
	}

	var timePlayed int
	var sessionTime int
	err = s.Server.db.QueryRow("SELECT time_played FROM characters WHERE id = $1", s.CharID).Scan(&timePlayed)
	if err != nil {
		panic(err)
	}

	sessionTime = int(Time_Current_Adjusted().Unix()) - int(s.SessionStart)
	timePlayed += sessionTime

	var rpGained int
	if s.Rights >= 0x40000000 { // N Course
		rpGained = timePlayed / 900
		timePlayed = timePlayed % 900
	} else {
		rpGained = timePlayed / 1800
		timePlayed = timePlayed % 1800
	}

	_, err = s.Server.db.Exec("UPDATE characters SET time_played = $1 WHERE id = $2", timePlayed, s.CharID)
	if err != nil {
		panic(err)
	}

	_, err = s.Server.db.Exec("UPDATE characters SET cafe_time=cafe_time+$1 WHERE id=$2", sessionTime, s.CharID)
	if err != nil {
		panic(err)
	}

	treasureHuntUnregister(s)

	if s.Stage == nil {
		return
	}

	s.Server.BroadcastMHF(&mhfpacket.MsgSysDeleteUser{
		CharID: s.CharID,
	}, s)

	s.Server.Lock()
	for _, stage := range s.Server.Stages {
		//nolint:gosimple
		if _, exists := stage.ReservedClientSlots[s.CharID]; exists {
			delete(stage.ReservedClientSlots, s.CharID)
		}
	}
	s.Server.Unlock()

	removeSessionFromSemaphore(s)
	removeSessionFromStage(s)

	saveData, err := GetCharacterSaveData(s, s.CharID)
	if err != nil {
		panic(err)
	}
	saveData.RP += uint16(rpGained)
	// nolint:ineffassign
	transaction, err := s.Server.db.Begin()
	if err != nil {
		// nolint:errcheck
		transaction.Rollback()
		panic(err)
	}

	err = saveData.Save(s, transaction)
	if err != nil {
		// nolint:errcheck // Error return value of `.` is not checked
		transaction.Rollback()
		panic(err)
	} else {
		// nolint:errcheck // Error return value of `.` is not checked
		transaction.Commit()
	}
}

func handleMsgSysSetStatus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysPing(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysPing)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgSysTime(s *Session, p mhfpacket.MHFPacket) {
	//pkt := p.(*mhfpacket.MsgSysTime)

	resp := &mhfpacket.MsgSysTime{
		GetRemoteTime: false,
		Timestamp:     uint32(Time_Current_Adjusted().Unix()), // JP timezone
	}
	s.QueueSendMHF(resp)
}

func handleMsgSysIssueLogkey(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysIssueLogkey)

	// Make a random log key for this session.
	logKey := make([]byte, 16)
	_, err := rand.Read(logKey)
	if err != nil {
		panic(err)
	}

	// TODO(Andoryuuta): In the offical client, the log key index is off by one,
	// cutting off the last byte in _most uses_. Find and document these accordingly.
	s.Lock()
	s.LogKey = logKey
	s.Unlock()

	// Issue it.
	resp := byteframe.NewByteFrame()
	resp.WriteBytes(logKey)
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgSysRecordLog(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysRecordLog)
	// remove a client returning to town from reserved slots to make sure the stage is hidden from board
	delete(s.Stage.ReservedClientSlots, s.CharID)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgSysEcho(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysLockGlobalSema(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysLockGlobalSema)

	bf := byteframe.NewByteFrame()
	// Unk
	// 0x00 when no ID sent
	// 0x02 when ID sent
	if pkt.ServerChannelIDLength == 1 {
		bf.WriteBytes([]byte{0x00, 0x00, 0x00, 0x01, 0x00})
	} else {
		bf.WriteUint8(0x02)
		bf.WriteUint8(0x00) // Unk
		bf.WriteUint16(uint16(pkt.ServerChannelIDLength))
		bf.WriteBytes([]byte(pkt.ServerChannelIDString))
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgSysUnlockGlobalSema(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysUnlockGlobalSema)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 8))
}

func handleMsgSysUpdateRight(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAuthQuery(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAuthTerminal(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysRightsReload(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysRightsReload)
	updateRights(s)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfTransitMessage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfTransitMessage)
	resp := byteframe.NewByteFrame()
	resp.WriteUint16(0)
	var count uint16
	switch pkt.SearchType {
	case 1: // CID
		bf := byteframe.NewByteFrameFromBytes(pkt.MessageData)
		CharID := bf.ReadUint32()
		for _, c := range s.Server.Channels {
			for _, session := range c.Sessions {
				if session.CharID == CharID {
					count++
					sessionName := stringsupport.UTF8ToSJIS(session.Name)
					sessionStage := stringsupport.UTF8ToSJIS(session.StageID)
					resp.WriteUint32(binary.LittleEndian.Uint32(net.ParseIP(c.IP).To4()))
					resp.WriteUint16(c.Port)
					resp.WriteUint32(session.CharID)
					resp.WriteBool(true)
					resp.WriteUint8(uint8(len(sessionName) + 1))
					resp.WriteUint16(uint16(len(c.userBinaryParts[userBinaryPartID{charID: session.CharID, index: 3}])))
					resp.WriteBytes(make([]byte, 40))
					resp.WriteUint8(uint8(len(sessionStage) + 1))
					resp.WriteBytes(make([]byte, 8))
					resp.WriteNullTerminatedBytes(sessionName)
					resp.WriteBytes(c.userBinaryParts[userBinaryPartID{session.CharID, 3}])
					resp.WriteNullTerminatedBytes(sessionStage)
				}
			}
		}
	case 2: // Name
		bf := byteframe.NewByteFrameFromBytes(pkt.MessageData)
		bf.ReadUint16() // lenSearchTerm
		bf.ReadUint16() // maxResults
		bf.ReadUint8()  // Unk
		searchTerm := stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes())
		for _, c := range s.Server.Channels {
			for _, session := range c.Sessions {
				if count == 100 {
					break
				}
				if strings.Contains(session.Name, searchTerm) {
					count++
					sessionName := stringsupport.UTF8ToSJIS(session.Name)
					sessionStage := stringsupport.UTF8ToSJIS(session.StageID)
					resp.WriteUint32(binary.LittleEndian.Uint32(net.ParseIP(c.IP).To4()))
					resp.WriteUint16(c.Port)
					resp.WriteUint32(session.CharID)
					resp.WriteBool(true)
					resp.WriteUint8(uint8(len(sessionName) + 1))
					resp.WriteUint16(uint16(len(c.userBinaryParts[userBinaryPartID{session.CharID, 3}])))
					resp.WriteBytes(make([]byte, 40))
					resp.WriteUint8(uint8(len(sessionStage) + 1))
					resp.WriteBytes(make([]byte, 8))
					resp.WriteNullTerminatedBytes(sessionName)
					resp.WriteBytes(c.userBinaryParts[userBinaryPartID{charID: session.CharID, index: 3}])
					resp.WriteNullTerminatedBytes(sessionStage)
				}
			}
		}
	case 3: // Enumerate Party
		bf := byteframe.NewByteFrameFromBytes(pkt.MessageData)
		ip := bf.ReadBytes(4)
		ipString := fmt.Sprintf("%d.%d.%d.%d", ip[3], ip[2], ip[1], ip[0])
		port := bf.ReadUint16()
		bf.ReadUint16() // lenStage
		maxResults := bf.ReadUint16()
		bf.ReadBytes(1)
		stageID := stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes())
		for _, c := range s.Server.Channels {
			if c.IP == ipString && c.Port == port {
				for _, stage := range c.Stages {
					if stage.Id == stageID {
						if count == maxResults {
							break
						}
						for session := range stage.Clients {
							count++
							hrp := uint16(1)
							gr := uint16(0)
							_ = s.Server.db.QueryRow("SELECT hrp, gr FROM characters WHERE id=$1", session.CharID).Scan(&hrp, &gr)
							sessionStage := stringsupport.UTF8ToSJIS(session.StageID)
							sessionName := stringsupport.UTF8ToSJIS(session.Name)
							resp.WriteUint32(binary.LittleEndian.Uint32(net.ParseIP(c.IP).To4()))
							resp.WriteUint16(c.Port)
							resp.WriteUint32(session.CharID)
							resp.WriteUint8(uint8(len(sessionStage) + 1))
							resp.WriteUint8(uint8(len(sessionName) + 1))
							resp.WriteUint8(0)
							resp.WriteUint8(7) // lenBinary
							resp.WriteBytes(make([]byte, 48))
							resp.WriteNullTerminatedBytes(sessionStage)
							resp.WriteNullTerminatedBytes(sessionName)
							resp.WriteUint16(hrp)
							resp.WriteUint16(gr)
							resp.WriteBytes([]byte{0x06, 0x10, 0x00}) // Unk
						}
					}
				}
			}
		}
	case 4: // Find Party
		bf := byteframe.NewByteFrameFromBytes(pkt.MessageData)
		setting := bf.ReadUint8()
		maxResults := bf.ReadUint16()
		// nolint:errcheck
		bf.Seek(2, 1)
		partyType := bf.ReadUint16()
		rankRestriction := uint16(0)
		if setting >= 2 {
			// nolint:errcheck
			bf.Seek(2, 1)
			rankRestriction = bf.ReadUint16()
		}
		targets := make([]uint16, 4)
		if setting >= 3 {
			// nolint:errcheck
			bf.Seek(1, 1)
			lenTargets := int(bf.ReadUint8())
			for i := 0; i < lenTargets; i++ {
				targets[i] = bf.ReadUint16()
			}
		}
		var stagePrefix string
		switch partyType {
		case 0: // Public Bar
			stagePrefix = "sl2Ls210"
		case 1: // Tokotoko Partnya
			stagePrefix = "sl2Ls463"
		case 2: // Hunting Prowess Match
			stagePrefix = "sl2Ls286"
		case 3: // Volpakkun Together
			stagePrefix = "sl2Ls465"
		case 5: // Quick Party
			// Unk
		}
		for _, c := range s.Server.Channels {
			for _, stage := range c.Stages {
				if count == maxResults {
					break
				}
				if strings.HasPrefix(stage.Id, stagePrefix) {
					sb3 := byteframe.NewByteFrameFromBytes(stage.RawBinaryData[StageBinaryKey{1, 3}])
					// nolint:errcheck
					sb3.Seek(4, 0)
					stageRankRestriction := sb3.ReadUint16()
					stageTarget := sb3.ReadUint16()
					if rankRestriction != 0xFFFF && stageRankRestriction < rankRestriction {
						continue
					}
					count++
					sessionStage := stringsupport.UTF8ToSJIS(stage.Id)
					resp.WriteUint32(binary.LittleEndian.Uint32(net.ParseIP(c.IP).To4()))
					resp.WriteUint16(c.Port)
					resp.WriteUint16(0) // Static?
					resp.WriteUint16(0) // Unk
					resp.WriteUint16(uint16(len(stage.Clients)))
					resp.WriteUint16(stage.MaxPlayers)
					resp.WriteUint16(0) // Num clients entered from stage
					resp.WriteUint16(stage.MaxPlayers)
					resp.WriteUint8(1) // Static?
					resp.WriteUint8(uint8(len(sessionStage) + 1))
					resp.WriteUint8(uint8(len(stage.RawBinaryData[StageBinaryKey{1, 0}])))
					resp.WriteUint8(uint8(len(stage.RawBinaryData[StageBinaryKey{1, 1}])))
					resp.WriteUint16(stageRankRestriction)
					resp.WriteUint16(stageTarget)
					resp.WriteBytes(make([]byte, 12))
					resp.WriteNullTerminatedBytes(sessionStage)
					resp.WriteBytes(stage.RawBinaryData[StageBinaryKey{1, 0}])
					resp.WriteBytes(stage.RawBinaryData[StageBinaryKey{1, 1}])
				}
			}
		}
	}
	if (pkt.SearchType == 1 || pkt.SearchType == 3) && count == 0 {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}
	// nolint:errcheck
	resp.Seek(0, io.SeekStart)
	resp.WriteUint16(count)
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgCaExchangeItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPresentBox(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfServerCommand(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAnnounce(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAnnounce)
	s.Server.BroadcastRaviente(pkt.IPAddress, pkt.Port, pkt.StageID, pkt.Type)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfSetLoginwindow(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysTransBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCollectBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysGetState(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSerialize(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEnumlobby(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysEnumuser(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysInfokyserver(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetCaUniqueID(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAcquireItem(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfTransferItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfTransferItem)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfEnumeratePrice(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumeratePrice)
	//resp := byteframe.NewByteFrame()
	//resp.WriteUint16(0) // Entry type 1 count
	//resp.WriteUint16(0) // Entry type 2 count
	// directly lifted for now because lacking it crashes the counter on having actual events present
	data, _ := hex.DecodeString("0000000066000003E800000000007300640100000320000000000006006401000003200000000000300064010000044C00000000007200640100000384000000000034006401000003840000000000140064010000051400000000006E006401000003E8000000000016006401000003E8000000000001006401000003200000000000430064010000057800000000006F006401000003840000000000330064010000044C00000000000B006401000003E800000000000F006401000006400000000000700064010000044C0000000000110064010000057800000000004C006401000003E8000000000059006401000006A400000000006D006401000005DC00000000004B006401000005DC000000000050006401000006400000000000350064010000070800000000006C0064010000044C000000000028006401000005DC00000000005300640100000640000000000060006401000005DC00000000005E0064010000051400000000007B006401000003E80000000000740064010000070800000000006B0064010000025800000000001B0064010000025800000000001C006401000002BC00000000001F006401000006A400000000007900640100000320000000000008006401000003E80000000000150064010000070800000000007A0064010000044C00000000000E00640100000640000000000055006401000007D0000000000002006401000005DC00000000002F0064010000064000000000002A0064010000076C00000000007E006401000002BC0000000000440064010000038400000000005C0064010000064000000000005B006401000006A400000000007D0064010000076C00000000007F006401000005DC0000000000540064010000064000000000002900640100000960000000000024006401000007D0000000000081006401000008340000000000800064010000038400000000001A006401000003E800000000002D0064010000038400000000004A006401000006A400000000005A00640100000384000000000027006401000007080000000000830064010000076C000000000040006401000006400000000000690064010000044C000000000025006401000004B000000000003100640100000708000000000082006401000003E800000000006500640100000640000000000051006401000007D000000000008C0064010000070800000000004D0064010000038400000000004E0064010000089800000000008B006401000004B000000000002E006401000009600000000000920064010000076C00000000008E00640100000514000000000068006401000004B000000000002B006401000003E800000000002C00640100000BB8000000000093006401000008FC00000000009000640100000AF0000000000094006401000006A400000000008D0064010000044C000000000052006401000005DC00000000004F006401000008980000000000970064010000070800000000006A0064010000064000000000005F00640100000384000000000026006401000008FC000000000096006401000007D00000000000980064010000076C000000000041006401000006A400000000003B006401000007080000000000360064010000083400000000009F00640100000A2800000000009A0064010000076C000000000021006401000007D000000000006300640100000A8C0000000000990064010000089800000000009E006401000007080000000000A100640100000C1C0000000000A200640100000C800000000000A400640100000DAC0000000000A600640100000C800000000000A50064010010")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfEnumerateOrder(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateOrder)
	stubEnumerateNoResults(s, pkt.AckHandle)
}

func handleMsgMhfGetExtraInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfEnumerateUnionItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateUnionItem)
	var boxContents []byte
	bf := byteframe.NewByteFrame()
	err := s.Server.db.QueryRow("SELECT item_box FROM users, characters WHERE characters.id = $1 AND users.id = characters.user_id", int(s.CharID)).Scan(&boxContents)
	if err != nil {
		s.logger.Fatal("Failed to get shared item box contents from db", zap.Error(err))
	} else {
		if len(boxContents) == 0 {
			bf.WriteUint32(0x00)
		} else {
			amount := len(boxContents) / 4
			bf.WriteUint16(uint16(amount))
			bf.WriteUint32(0x00)
			bf.WriteUint16(0x00)
			for i := 0; i < amount; i++ {
				bf.WriteUint32(binary.BigEndian.Uint32(boxContents[i*4 : i*4+4]))
				if i+1 != amount {
					bf.WriteUint64(0x00)
				}
			}
		}
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfUpdateUnionItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateUnionItem)
	// Get item cache from DB
	var boxContents []byte
	var oldItems []Item

	err := s.Server.db.QueryRow("SELECT item_box FROM users, characters WHERE characters.id = $1 AND users.id = characters.user_id", int(s.CharID)).Scan(&boxContents)
	if err != nil {
		s.logger.Fatal("Failed to get shared item box contents from db", zap.Error(err))
	} else {
		amount := len(boxContents) / 4
		oldItems = make([]Item, amount)
		for i := 0; i < amount; i++ {
			oldItems[i].ItemId = binary.BigEndian.Uint16(boxContents[i*4 : i*4+2])
			oldItems[i].Amount = binary.BigEndian.Uint16(boxContents[i*4+2 : i*4+4])
		}
	}

	// Update item stacks
	newItems := make([]Item, len(oldItems))
	copy(newItems, oldItems)
	for i := 0; i < int(pkt.Amount); i++ {
		for j := 0; j <= len(oldItems); j++ {
			if j == len(oldItems) {
				var newItem Item
				newItem.ItemId = pkt.Items[i].ItemId
				newItem.Amount = pkt.Items[i].Amount
				newItems = append(newItems, newItem)
				break
			}
			if pkt.Items[i].ItemId == oldItems[j].ItemId {
				newItems[j].Amount = pkt.Items[i].Amount
				break
			}
		}
	}

	// Delete empty item stacks
	for i := len(newItems) - 1; i >= 0; i-- {
		if int(newItems[i].Amount) == 0 {
			copy(newItems[i:], newItems[i+1:])
			newItems[len(newItems)-1] = make([]Item, 1)[0]
			newItems = newItems[:len(newItems)-1]
		}
	}

	// Create new item cache
	bf := byteframe.NewByteFrame()
	for i := 0; i < len(newItems); i++ {
		bf.WriteUint16(newItems[i].ItemId)
		bf.WriteUint16(newItems[i].Amount)
	}

	// Upload new item cache
	_, err = s.Server.db.Exec("UPDATE users SET item_box = $1 FROM characters WHERE  users.id = characters.user_id AND characters.id = $2", bf.Data(), int(s.CharID))
	if err != nil {
		s.logger.Fatal("Failed to update shared item box contents in db", zap.Error(err))
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfGetCogInfo(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfCheckWeeklyStamp(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCheckWeeklyStamp)

	if s.Server.erupeConfig.DevModeOptions.DisableStamps {
		resp := byteframe.NewByteFrame()
		resp.WriteUint16(0x000E)
		resp.WriteUint16(0x0001)
		resp.WriteUint16(0x0000)
		resp.WriteUint16(0x0000) // 0x0000 stops the vaguely annoying log in pop up
		resp.WriteUint32(0)
		resp.WriteUint32(0x5dddcbb3) // Timestamp

		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
		return
	}

	weekCurrentStart := TimeWeekStart()
	weekNextStart := TimeWeekNext()
	var total, redeemed, updated uint16
	var nextClaim time.Time
	err := s.Server.db.QueryRow(fmt.Sprintf("SELECT %s_next FROM stamps WHERE character_id=$1", pkt.StampType), s.CharID).Scan(&nextClaim)
	if err != nil {
		_, err = s.Server.db.Exec("INSERT INTO stamps (character_id, hl_next, ex_next) VALUES ($1, $2, $2)", s.CharID, weekNextStart)
		if err != nil {
			doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
			return
		}

		nextClaim = weekNextStart
	}
	if nextClaim.Before(weekCurrentStart) {
		_, err := s.Server.db.Exec(fmt.Sprintf("UPDATE stamps SET %s_total=%s_total+1, %s_next=$1 WHERE character_id=$2", pkt.StampType, pkt.StampType, pkt.StampType), weekNextStart, s.CharID)
		if err != nil {
			s.logger.Error("INVALID UPDATE STAMPS ON handleMsgMhfCheckWeeklyStamp", zap.Error(err))
			doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
			return
		}

		updated = 1
	}
	err = s.Server.db.QueryRow(fmt.Sprintf("SELECT %s_total, %s_redeemed FROM stamps WHERE character_id=$1", pkt.StampType, pkt.StampType), s.CharID).Scan(&total, &redeemed)
	if err != nil {
		s.logger.Error("FAILED TO SELECT STAMPS ON handleMsgMhfCheckWeeklyStamp", zap.Error(err))
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	bf := byteframe.NewByteFrame()
	bf.WriteUint16(total)
	bf.WriteUint16(redeemed)
	bf.WriteUint16(updated)
	bf.WriteUint32(0) // Unk
	bf.WriteUint32(uint32(weekCurrentStart.Unix()))
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfExchangeWeeklyStamp(s *Session, p mhfpacket.MHFPacket) {
	if s.Server.erupeConfig.DevModeOptions.DisableStamps {
		return
	}

	pkt := p.(*mhfpacket.MsgMhfExchangeWeeklyStamp)
	var total, redeemed uint16
	var tktStack mhfpacket.WarehouseStack
	if pkt.Unk1 == 0xA { // Yearly Sub Ex
		err := s.Server.db.QueryRow("UPDATE stamps SET hl_total=hl_total-48, hl_redeemed=hl_redeemed-48 WHERE character_id=$1 RETURNING hl_total, hl_redeemed", s.CharID).Scan(&total, &redeemed)
		if err != nil {
			doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
			return
		}

		tktStack = mhfpacket.WarehouseStack{ItemID: 0x08A2, Quantity: 1}
	} else {
		err := s.Server.db.QueryRow(fmt.Sprintf("UPDATE stamps SET %s_redeemed=%s_redeemed+8 WHERE character_id=$1 RETURNING %s_total, %s_redeemed", pkt.StampType, pkt.StampType, pkt.StampType, pkt.StampType), s.CharID).Scan(&total, &redeemed)
		if err != nil {
			doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
			return
		}

		if pkt.StampType == "hl" {
			tktStack = mhfpacket.WarehouseStack{ItemID: 0x065E, Quantity: 5}
		} else {
			tktStack = mhfpacket.WarehouseStack{ItemID: 0x065F, Quantity: 5}
		}
	}
	addWarehouseGift(s, "item", tktStack)
	bf := byteframe.NewByteFrame()
	bf.WriteUint16(total)
	bf.WriteUint16(redeemed)
	bf.WriteUint16(0)
	bf.WriteUint32(0) // Unk, but has possible values
	bf.WriteUint32(uint32(TimeWeekStart().Unix()))
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func getGookData(s *Session, cid uint32) (uint16, []byte) {
	var data []byte
	var count uint16
	bf := byteframe.NewByteFrame()
	for i := 0; i < 5; i++ {
		err := s.Server.db.QueryRow(fmt.Sprintf("SELECT gook%d FROM gook WHERE id=$1", i), cid).Scan(&data)
		if err != nil {
			// nolint:errcheck // Error return value of `.` is not checked
			s.Server.db.Exec("INSERT INTO gook (id) VALUES ($1)", s.CharID)
			return 0, bf.Data()
		}
		if err == nil && data != nil {
			count++
			if s.CharID == cid && count == 1 {
				gook := byteframe.NewByteFrameFromBytes(data)
				bf.WriteBytes(gook.ReadBytes(4))
				d := gook.ReadBytes(2)
				bf.WriteBytes(d)
				bf.WriteBytes(d)
				bf.WriteBytes(gook.DataFromCurrent())
			} else {
				bf.WriteBytes(data)
			}
		}
	}
	return count, bf.Data()
}

func handleMsgMhfEnumerateGuacot(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateGuacot)
	bf := byteframe.NewByteFrame()
	count, data := getGookData(s, s.CharID)
	bf.WriteUint16(count)
	bf.WriteBytes(data)
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfUpdateGuacot(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateGuacot)
	for _, gook := range pkt.Gooks {
		if !gook.Exists {
			// nolint:errcheck
			s.Server.db.Exec(fmt.Sprintf("UPDATE gook SET gook%d=NULL WHERE id=$1", gook.Index), s.CharID)
		} else {
			bf := byteframe.NewByteFrame()
			bf.WriteUint32(gook.Index)
			bf.WriteUint16(gook.Type)
			bf.WriteBytes(gook.Data)
			bf.WriteUint8(gook.NameLen)
			bf.WriteBytes(gook.Name)
			// nolint:errcheck // Error return value of `.` is not checked
			s.Server.db.Exec(fmt.Sprintf("UPDATE gook SET gook%d=$1 WHERE id=$2", gook.Index), bf.Data(), s.CharID)
		}
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfInfoScenarioCounter(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfInfoScenarioCounter)
	scenarioCounter := []struct {
		MainID uint32
		Unk1   uint8 // Bool item exchange?
		// 0 = basic, 1 = veteran, 3 = other, 6 = pallone, 7 = diva
		CategoryID uint8
	}{
		//000000110000
		{
			MainID: 0x00000011, Unk1: 0, CategoryID: 0,
		},
		// 0000005D0001
		{
			MainID: 0x0000005D, Unk1: 0, CategoryID: 1,
		},
		// 0000005C0001
		{
			MainID: 0x0000005C, Unk1: 0, CategoryID: 1,
		},
		// 000000510001
		{
			MainID: 0x00000051, Unk1: 0, CategoryID: 1,
		},
		// 0000005B0001
		{
			MainID: 0x0000005B, Unk1: 0, CategoryID: 1,
		},
		// 0000005A0001
		{
			MainID: 0x0000005A, Unk1: 0, CategoryID: 1,
		},
		// 000000590001
		{
			MainID: 0x00000059, Unk1: 0, CategoryID: 1,
		},
		// 000000580001
		{
			MainID: 0x00000058, Unk1: 0, CategoryID: 1,
		},
		// 000000570001
		{
			MainID: 0x00000057, Unk1: 0, CategoryID: 1,
		},
		// 000000560001
		{
			MainID: 0x00000056, Unk1: 0, CategoryID: 1,
		},
		// 000000550001
		{
			MainID: 0x00000055, Unk1: 0, CategoryID: 1,
		},
		// 000000540001
		{
			MainID: 0x00000054, Unk1: 0, CategoryID: 1,
		},
		// 000000530001
		{
			MainID: 0x00000053, Unk1: 0, CategoryID: 1,
		},
		// 000000520001
		{
			MainID: 0x00000052, Unk1: 0, CategoryID: 1,
		},
		// 000000570103
		{
			MainID: 0x00000057, Unk1: 1, CategoryID: 3,
		},
		// 000000580103
		{
			MainID: 0x00000058, Unk1: 1, CategoryID: 3,
		},
		// 000000590103
		{
			MainID: 0x00000059, Unk1: 1, CategoryID: 3,
		},
		// 0000005A0103
		{
			MainID: 0x0000005A, Unk1: 1, CategoryID: 3,
		},
		// 0000005B0103
		{
			MainID: 0x0000005B, Unk1: 1, CategoryID: 3,
		},
		// 0000005C0103
		{
			MainID: 0x0000005C, Unk1: 1, CategoryID: 3,
		},
		// 000000530103
		{
			MainID: 0x00000053, Unk1: 1, CategoryID: 3,
		},
		// 000000560103
		{
			MainID: 0x00000056, Unk1: 1, CategoryID: 3,
		},
		// 0000003C0103
		{
			MainID: 0x0000003C, Unk1: 1, CategoryID: 3,
		},
		// 0000003A0103
		{
			MainID: 0x0000003A, Unk1: 1, CategoryID: 3,
		},
		// 0000003B0103
		{
			MainID: 0x0000003B, Unk1: 1, CategoryID: 3,
		},
		// 0000001B0103
		{
			MainID: 0x0000001B, Unk1: 1, CategoryID: 3,
		},
		// 000000190103
		{
			MainID: 0x00000019, Unk1: 1, CategoryID: 3,
		},
		// 0000001A0103
		{
			MainID: 0x0000001A, Unk1: 1, CategoryID: 3,
		},
		// 000000170103
		{
			MainID: 0x00000017, Unk1: 1, CategoryID: 3,
		},
		// 000000020103
		{
			MainID: 0x00000002, Unk1: 1, CategoryID: 3,
		},
		// 000000030103
		{
			MainID: 0x00000003, Unk1: 1, CategoryID: 3,
		},
		// 000000040103
		{
			MainID: 0x00000004, Unk1: 1, CategoryID: 3,
		},
		// 0000001F0103
		{
			MainID: 0x0000001F, Unk1: 1, CategoryID: 3,
		},
		// 000000200103
		{
			MainID: 0x00000020, Unk1: 1, CategoryID: 3,
		},
		// 000000210103
		{
			MainID: 0x00000021, Unk1: 1, CategoryID: 3,
		},
		// 000000220103
		{
			MainID: 0x00000022, Unk1: 1, CategoryID: 3,
		},
		// 000000230103
		{
			MainID: 0x00000023, Unk1: 1, CategoryID: 3,
		},
		// 000000240103
		{
			MainID: 0x00000024, Unk1: 1, CategoryID: 3,
		},
		// 000000250103
		{
			MainID: 0x00000025, Unk1: 1, CategoryID: 3,
		},
		// 000000280103
		{
			MainID: 0x00000028, Unk1: 1, CategoryID: 3,
		},
		// 000000260103
		{
			MainID: 0x00000026, Unk1: 1, CategoryID: 3,
		},
		// 000000270103
		{
			MainID: 0x00000027, Unk1: 1, CategoryID: 3,
		},
		// 000000300103
		{
			MainID: 0x00000030, Unk1: 1, CategoryID: 3,
		},
		// 0000000C0103
		{
			MainID: 0x0000000C, Unk1: 1, CategoryID: 3,
		},
		// 0000000D0103
		{
			MainID: 0x0000000D, Unk1: 1, CategoryID: 3,
		},
		// 0000001E0103
		{
			MainID: 0x0000001E, Unk1: 1, CategoryID: 3,
		},
		// 0000001D0103
		{
			MainID: 0x0000001D, Unk1: 1, CategoryID: 3,
		},
		// 0000002E0003
		{
			MainID: 0x0000002E, Unk1: 0, CategoryID: 3,
		},
		// 000000000004
		{
			MainID: 0x00000000, Unk1: 0, CategoryID: 4,
		},
		// 000000010004
		{
			MainID: 0x00000001, Unk1: 0, CategoryID: 4,
		},
		// 000000020004
		{
			MainID: 0x00000002, Unk1: 0, CategoryID: 4,
		},
		// 000000030004
		{
			MainID: 0x00000003, Unk1: 0, CategoryID: 4,
		},
		// 000000040004
		{
			MainID: 0x00000004, Unk1: 0, CategoryID: 4,
		},
		// 000000050004
		{
			MainID: 0x00000005, Unk1: 0, CategoryID: 4,
		},
		// 000000060004
		{
			MainID: 0x00000006, Unk1: 0, CategoryID: 4,
		},
		// 000000070004
		{
			MainID: 0x00000007, Unk1: 0, CategoryID: 4,
		},
		// 000000080004
		{
			MainID: 0x00000008, Unk1: 0, CategoryID: 4,
		},
		// 000000090004
		{
			MainID: 0x00000009, Unk1: 0, CategoryID: 4,
		},
		// 0000000A0004
		{
			MainID: 0x0000000A, Unk1: 0, CategoryID: 4,
		},
		// 0000000B0004
		{
			MainID: 0x0000000B, Unk1: 0, CategoryID: 4,
		},
		// 0000000C0004
		{
			MainID: 0x0000000C, Unk1: 0, CategoryID: 4,
		},
		// 0000000D0004
		{
			MainID: 0x0000000D, Unk1: 0, CategoryID: 4,
		},
		// 0000000E0004
		{
			MainID: 0x0000000E, Unk1: 0, CategoryID: 4,
		},
		// 000000320005
		{
			MainID: 0x00000032, Unk1: 0, CategoryID: 5,
		},
		// 000000330005
		{
			MainID: 0x00000033, Unk1: 0, CategoryID: 5,
		},
		// 000000340005
		{
			MainID: 0x00000034, Unk1: 0, CategoryID: 5,
		},
		// 000000350005
		{
			MainID: 0x00000035, Unk1: 0, CategoryID: 5,
		},
		// 000000360005
		{
			MainID: 0x00000036, Unk1: 0, CategoryID: 5,
		},
		// 000000370005
		{
			MainID: 0x00000037, Unk1: 0, CategoryID: 5,
		},
		// 000000380005
		{
			MainID: 0x00000038, Unk1: 0, CategoryID: 5,
		},
		// 0000003A0005
		{
			MainID: 0x0000003A, Unk1: 0, CategoryID: 5,
		},
		// 0000003F0005
		{
			MainID: 0x0000003F, Unk1: 0, CategoryID: 5,
		},
		// 000000400005
		{
			MainID: 0x00000040, Unk1: 0, CategoryID: 5,
		},
		// 000000410005
		{
			MainID: 0x00000041, Unk1: 0, CategoryID: 5,
		},
		// 000000430005
		{
			MainID: 0x00000043, Unk1: 0, CategoryID: 5,
		},
		// 000000470005
		{
			MainID: 0x00000047, Unk1: 0, CategoryID: 5,
		},
		// 0000004B0005
		{
			MainID: 0x0000004B, Unk1: 0, CategoryID: 5,
		},
		// 0000003D0005
		{
			MainID: 0x0000003D, Unk1: 0, CategoryID: 5,
		},
		// 000000440005
		{
			MainID: 0x00000044, Unk1: 0, CategoryID: 5,
		},
		// 000000420005
		{
			MainID: 0x00000042, Unk1: 0, CategoryID: 5,
		},
		// 0000004C0005
		{
			MainID: 0x0000004C, Unk1: 0, CategoryID: 5,
		},
		// 000000460005
		{
			MainID: 0x00000046, Unk1: 0, CategoryID: 5,
		},
		// 0000004D0005
		{
			MainID: 0x0000004D, Unk1: 0, CategoryID: 5,
		},
		// 000000480005
		{
			MainID: 0x00000048, Unk1: 0, CategoryID: 5,
		},
		// 0000004A0005
		{
			MainID: 0x0000004A, Unk1: 0, CategoryID: 5,
		},
		// 000000490005
		{
			MainID: 0x00000049, Unk1: 0, CategoryID: 5,
		},
		// 0000004E0005
		{
			MainID: 0x0000004E, Unk1: 0, CategoryID: 5,
		},
		// 000000450005
		{
			MainID: 0x00000045, Unk1: 0, CategoryID: 5,
		},
		// 0000003E0005
		{
			MainID: 0x0000003E, Unk1: 0, CategoryID: 5,
		},
		// 0000004F0005
		{
			MainID: 0x0000004F, Unk1: 0, CategoryID: 5,
		},
		// 000000000106
		{
			MainID: 0x00000000, Unk1: 1, CategoryID: 6,
		},
		// 000000010106
		{
			MainID: 0x00000001, Unk1: 1, CategoryID: 6,
		},
		// 000000020106
		{
			MainID: 0x00000002, Unk1: 1, CategoryID: 6,
		},
		// 000000030106
		{
			MainID: 0x00000003, Unk1: 1, CategoryID: 6,
		},
		// 000000040106
		{
			MainID: 0x00000004, Unk1: 1, CategoryID: 6,
		},
		// 000000050106
		{
			MainID: 0x00000005, Unk1: 1, CategoryID: 6,
		},
		// 000000060106
		{
			MainID: 0x00000006, Unk1: 1, CategoryID: 6,
		},
		// 000000070106
		{
			MainID: 0x00000007, Unk1: 1, CategoryID: 6,
		},
		// 000000080106
		{
			MainID: 0x00000008, Unk1: 1, CategoryID: 6,
		},
		// 000000090106
		{
			MainID: 0x00000009, Unk1: 1, CategoryID: 6,
		},
		// 000000110106
		{
			MainID: 0x00000011, Unk1: 1, CategoryID: 6,
		},
		// 0000000A0106
		{
			MainID: 0x0000000A, Unk1: 1, CategoryID: 6,
		},
		// 0000000B0106
		{
			MainID: 0x0000000B, Unk1: 1, CategoryID: 6,
		},
		// 0000000C0106
		{
			MainID: 0x0000000C, Unk1: 1, CategoryID: 6,
		},
		// 0000000D0106
		{
			MainID: 0x0000000D, Unk1: 1, CategoryID: 6,
		},
		// 0000000E0106
		{
			MainID: 0x0000000E, Unk1: 1, CategoryID: 6,
		},
		// 0000000F0106
		{
			MainID: 0x0000000F, Unk1: 1, CategoryID: 6,
		},
		// 000000100106
		{
			MainID: 0x00000010, Unk1: 1, CategoryID: 6,
		},
		// 000000320107
		{
			MainID: 0x00000032, Unk1: 1, CategoryID: 7,
		},
		// 000000350107
		{
			MainID: 0x00000035, Unk1: 1, CategoryID: 7,
		},
		// 0000003E0107
		{
			MainID: 0x0000003E, Unk1: 1, CategoryID: 7,
		},
		// 000000340107
		{
			MainID: 0x00000034, Unk1: 1, CategoryID: 7,
		},
		// 000000380107
		{
			MainID: 0x00000038, Unk1: 1, CategoryID: 7,
		},
		// 000000330107
		{
			MainID: 0x00000033, Unk1: 1, CategoryID: 7,
		},
		// 000000310107
		{
			MainID: 0x00000031, Unk1: 1, CategoryID: 7,
		},
		// 000000360107
		{
			MainID: 0x00000036, Unk1: 1, CategoryID: 7,
		},
		// 000000390107
		{
			MainID: 0x00000039, Unk1: 1, CategoryID: 7,
		},
		// 000000370107
		{
			MainID: 0x00000037, Unk1: 1, CategoryID: 7,
		},
		// 0000003D0107
		{
			MainID: 0x0000003D, Unk1: 1, CategoryID: 7,
		},
		// 0000003A0107
		{
			MainID: 0x0000003A, Unk1: 1, CategoryID: 7,
		},
		// 0000003C0107
		{
			MainID: 0x0000003C, Unk1: 1, CategoryID: 7,
		},
		// 0000003B0107
		{
			MainID: 0x0000003B, Unk1: 1, CategoryID: 7,
		},
		// 0000002A0107
		{
			MainID: 0x0000002A, Unk1: 1, CategoryID: 7,
		},
		// 000000300107
		{
			MainID: 0x00000030, Unk1: 1, CategoryID: 7,
		},
		// 000000280107
		{
			MainID: 0x00000028, Unk1: 1, CategoryID: 7,
		},
		// 000000270107
		{
			MainID: 0x00000027, Unk1: 1, CategoryID: 7,
		},
		// 0000002B0107
		{
			MainID: 0x0000002B, Unk1: 1, CategoryID: 7,
		},
		// 0000002E0107
		{
			MainID: 0x0000002E, Unk1: 1, CategoryID: 7,
		},
		// 000000290107
		{
			MainID: 0x00000029, Unk1: 1, CategoryID: 7,
		},
		// 0000002C0107
		{
			MainID: 0x0000002C, Unk1: 1, CategoryID: 7,
		},
		// 0000002D0107
		{
			MainID: 0x0000002D, Unk1: 1, CategoryID: 7,
		},
		// 0000002F0107
		{
			MainID: 0x0000002F, Unk1: 1, CategoryID: 7,
		},
		// 000000250107
		{
			MainID: 0x00000025, Unk1: 1, CategoryID: 7,
		},
		// 000000220107
		{
			MainID: 0x00000022, Unk1: 1, CategoryID: 7,
		},
		// 000000210107
		{
			MainID: 0x00000021, Unk1: 1, CategoryID: 7,
		},
		// 000000200107
		{
			MainID: 0x00000020, Unk1: 1, CategoryID: 7,
		},
		// 0000001C0107
		{
			MainID: 0x0000001C, Unk1: 1, CategoryID: 7,
		},
		// 0000001A0107
		{
			MainID: 0x0000001A, Unk1: 1, CategoryID: 7,
		},
		// 000000240107
		{
			MainID: 0x00000024, Unk1: 1, CategoryID: 7,
		},
		// 000000260107
		{
			MainID: 0x00000026, Unk1: 1, CategoryID: 7,
		},
		// 000000230107
		{
			MainID: 0x00000023, Unk1: 1, CategoryID: 7,
		},
		// 0000001B0107
		{
			MainID: 0x0000001B, Unk1: 1, CategoryID: 7,
		},
		// 0000001E0107
		{
			MainID: 0x0000001E, Unk1: 1, CategoryID: 7,
		},
		// 0000001F0107
		{
			MainID: 0x0000001F, Unk1: 1, CategoryID: 7,
		},
		// 0000001D0107
		{
			MainID: 0x0000001D, Unk1: 1, CategoryID: 7,
		},
		// 000000180107
		{
			MainID: 0x00000018, Unk1: 1, CategoryID: 7,
		},
		// 000000170107
		{
			MainID: 0x00000017, Unk1: 1, CategoryID: 7,
		},
		// 000000160107
		{
			MainID: 0x00000016, Unk1: 1, CategoryID: 7,
		},
		// 000000150107
		// Missing file
		// {
		// 	MainID: 0x00000015, Unk1: 1, CategoryID: 7,
		// },
		// 000000190107
		{
			MainID: 0x00000019, Unk1: 1, CategoryID: 7,
		},
		// 000000140107
		// Missing file
		// {
		// 	MainID: 0x00000014, Unk1: 1, CategoryID: 7,
		// },
		// 000000070107
		// Missing file
		// {
		//	MainID: 0x00000007, Unk1: 1, CategoryID: 7,
		// },
		// 000000090107
		// Missing file
		// {
		//	MainID: 0x00000009, Unk1: 1, CategoryID: 7,
		// },
		// 0000000D0107
		// Missing file
		// {
		//	MainID: 0x0000000D, Unk1: 1, CategoryID: 7,
		// },
		// 000000100107
		// Missing file
		// {
		//	MainID: 0x00000010, Unk1: 1, CategoryID: 7,
		// },
		// 0000000C0107
		// Missing file
		// {
		//	MainID: 0x0000000C, Unk1: 1, CategoryID: 7,
		// },
		// 0000000E0107
		// Missing file
		// {
		//	MainID: 0x0000000E, Unk1: 1, CategoryID: 7,
		// },
		// 0000000F0107
		// Missing file
		// {
		//	MainID: 0x0000000F, Unk1: 1, CategoryID: 7,
		// },
		// 000000130107
		// Missing file
		// {
		//	MainID: 0x00000013, Unk1: 1, CategoryID: 7,
		// },
		// 0000000A0107
		// Missing file
		// {
		//	MainID: 0x0000000A, Unk1: 1, CategoryID: 7,
		// },
		// 000000080107
		// Missing file
		// {
		//	MainID: 0x00000008, Unk1: 1, CategoryID: 7,
		// },
		// 0000000B0107
		// Missing file
		// {
		//	MainID: 0x0000000B, Unk1: 1, CategoryID: 7,
		// },
		// 000000120107
		// Missing file
		// {
		//	MainID: 0x00000012, Unk1: 1, CategoryID: 7,
		// },
		// 000000110107
		// Missing file
		// {
		// 	MainID: 0x00000011, Unk1: 1, CategoryID: 7,
		// },
		// 000000060107
		// Missing file
		// {
		// 	MainID: 0x00000006, Unk1: 1, CategoryID: 7,
		// },
		// 000000050107
		// Missing file
		// {
		// 	MainID: 0x00000005, Unk1: 1, CategoryID: 7,
		// },
		// 000000040107
		// Missing file
		// {
		//	MainID: 0x00000004, Unk1: 1, CategoryID: 7,
		// },
		// 000000030107
		{
			MainID: 0x00000003, Unk1: 1, CategoryID: 7,
		},
		// 000000020107
		{
			MainID: 0x00000002, Unk1: 1, CategoryID: 7,
		},
		// 000000010107
		{
			MainID: 0x00000001, Unk1: 1, CategoryID: 7,
		},
		// 000000000107
		{
			MainID: 0x00000000, Unk1: 1, CategoryID: 7,
		},
	}

	resp := byteframe.NewByteFrame()
	resp.WriteUint8(uint8(len(scenarioCounter))) // Entry count
	for _, entry := range scenarioCounter {
		resp.WriteUint32(entry.MainID)
		resp.WriteUint8(entry.Unk1)
		resp.WriteUint8(entry.CategoryID)
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfGetBbsSnsStatus(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfApplyBbsArticle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEtcPoints(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetEtcPoints)

	resp := byteframe.NewByteFrame()
	resp.WriteUint8(0x3) // Maybe a count of uint32(s)?
	resp.WriteUint32(0)
	resp.WriteUint32(14)
	resp.WriteUint32(14)

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfUpdateEtcPoint(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateEtcPoint)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfStampcardStamp(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfStampcardStamp)
	// TODO: Work out where it gets existing stamp count from, its format and then
	// update the actual sent values to be correct
	doAckBufSucceed(s, pkt.AckHandle, []byte{0x03, 0xe7, 0x03, 0xe7, 0x02, 0x99, 0x02, 0x9c, 0x00, 0x00, 0x00, 0x00, 0x14, 0xf8, 0x69, 0x54})
}

func handleMsgMhfStampcardPrize(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUnreserveSrg(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadBeatLevel(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadBeatLevel)

	// This response is fixed and will never change on JP,
	// but I've left it dynamic for possible other client differences.
	resp := byteframe.NewByteFrame()
	for i := 0; i < int(pkt.ValidIDCount); i++ {
		resp.WriteUint32(pkt.IDs[i])
		resp.WriteUint32(1)
		resp.WriteUint32(1)
		resp.WriteUint32(1)
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfUpdateBeatLevel(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateBeatLevel)

	doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfReadBeatLevelAllRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadBeatLevelMyRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfReadLastWeekBeatRanking(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetFixedSeibatuRankingTable(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfKickExportForce(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEarthStatus(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetEarthStatus)

	// TODO(Andoryuuta): Track down format for this data,
	//	it can somehow be parsed as 8*uint32 chunks if the header is right.
	/*
		BEFORE ack-refactor:
			resp := byteframe.NewByteFrame()
			resp.WriteUint32(0)
			resp.WriteUint32(0)

			s.QueueAck(pkt.AckHandle, resp.Data())
	*/
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfRegistSpabiTime(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEarthValue(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetEarthValue)
	var earthValues []struct{ Unk0, Unk1, Unk2, Unk3, Unk4, Unk5 uint32 }
	if pkt.ReqType == 3 {
		earthValues = []struct {
			Unk0, Unk1, Unk2, Unk3, Unk4, Unk5 uint32
		}{
			// TW identical to JP
			{
				Unk0: 0x03E9,
				Unk1: 0x24,
			},
			{
				Unk0: 0x2329,
				Unk1: 0x03,
			},
			{
				Unk0: 0x232A,
				Unk1: 0x0A,
				Unk2: 0x012C,
			},
		}
	} else if pkt.ReqType == 2 {
		earthValues = []struct {
			Unk0, Unk1, Unk2, Unk3, Unk4, Unk5 uint32
		}{
			// JP response was empty
			{
				Unk0: 0x01,
				Unk1: 0x168B,
			},
			{
				Unk0: 0x02,
				Unk1: 0x0737,
			},
		}
	} else if pkt.ReqType == 1 {
		earthValues = []struct {
			Unk0, Unk1, Unk2, Unk3, Unk4, Unk5 uint32
		}{
			// JP simply sent 01 and 02 respectively
			{
				Unk0: 0x01,
				Unk1: 0x0138,
			},
			{
				Unk0: 0x02,
				Unk1: 0x63,
			},
		}
	}

	resp := byteframe.NewByteFrame()
	resp.WriteUint32(0x0A218EAD)               // Unk shared ID. Sent in response of MSG_MHF_GET_TOWER_INFO, MSG_MHF_GET_PAPER_DATA etc.
	resp.WriteUint32(0)                        // Unk
	resp.WriteUint32(0)                        // Unk
	resp.WriteUint32(uint32(len(earthValues))) // value count
	for _, v := range earthValues {
		resp.WriteUint32(v.Unk0)
		resp.WriteUint32(v.Unk1)
		resp.WriteUint32(v.Unk2)
		resp.WriteUint32(v.Unk3)
		resp.WriteUint32(v.Unk4)
		resp.WriteUint32(v.Unk5)
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfDebugPostValue(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetNotice(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostNotice(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRandFromTable(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetTinyBin(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetTinyBin)
	// requested after conquest quests
	doAckBufSucceed(s, pkt.AckHandle, []byte{})
}

func handleMsgMhfPostTinyBin(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetSenyuDailyCount(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetSeibattle(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetSeibattle)
	stubGetNoResults(s, pkt.AckHandle)
}

func handleMsgMhfPostSeibattle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRyoudama(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetRyoudama)
	// likely guild related
	// REQ: 00 04 13 53 8F 18 00
	// RSP: 0A 21 8E AD 00 00 00 00 00 00 00 00 00 00 00 01 00 01 FE 4E
	// REQ: 00 06 13 53 8F 18 00
	// RSP: 0A 21 8E AD 00 00 00 00 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00 00
	// REQ: 00 05 13 53 8F 18 00
	// RSP: 0A 21 8E AD 00 00 00 00 00 00 00 00 00 00 00 0E 2A 15 9E CC 00 00 00 01 82 79 83 4E 83 8A 81 5B 83 69 00 00 00 00 1E 55 B0 2F 00 00 00 01 8D F7 00 00 00 00 00 00 00 00 00 00 00 00 2A 15 9E CC 00 00 00 02 82 79 83 4E 83 8A 81 5B 83 69 00 00 00 00 03 D5 30 56 00 00 00 02 95 BD 91 F2 97 42 00 00 00 00 00 00 00 00 3F 57 76 9F 00 00 00 03 93 56 92 6E 96 B3 97 70 00 00 00 00 00 00 38 D9 0E C4 00 00 00 03 87 64 83 78 83 42 00 00 00 00 00 00 00 00 23 F3 B9 77 00 00 00 04 82 B3 82 CC 82 DC 82 E9 81 99 00 00 00 00 3F 1B 17 9C 00 00 00 04 82 B1 82 A4 82 BD 00 00 00 00 00 00 00 00 00 B9 F9 C0 00 00 00 05 82 CD 82 E9 82 A9 00 00 00 00 00 00 00 00 23 9F 9A EA 00 00 00 05 83 70 83 62 83 4C 83 83 83 49 00 00 00 00 38 D9 0E C4 00 00 00 06 87 64 83 78 83 42 00 00 00 00 00 00 00 00 1E 55 B0 2F 00 00 00 06 8D F7 00 00 00 00 00 00 00 00 00 00 00 00 03 D5 30 56 00 00 00 07 95 BD 91 F2 97 42 00 00 00 00 00 00 00 00 02 D3 B8 77 00 00 00 07 6F 77 6C 32 35 32 35 00 00 00 00 00 00 00
	data, _ := hex.DecodeString("0A218EAD0000000000000000000000010000000000000000")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfPostRyoudama(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetTenrouirai(s *Session, p mhfpacket.MHFPacket) {
	// if the game gets bad responses for this it breaks the ability to save
	pkt := p.(*mhfpacket.MsgMhfGetTenrouirai)
	var data []byte
	var err error
	if pkt.Unk0 == 1 {
		data, err = hex.DecodeString("0A218EAD000000000000000000000001010000000000060010")
	} else if pkt.Unk2 == 4 {
		data, err = hex.DecodeString("0A218EAD0000000000000000000000210101005000000202010102020104001000000202010102020106003200000202010002020104000C003202020101020201030032000002020101020202059C4000000202010002020105C35000320202010102020201003C00000202010102020203003200000201010001020203002800320201010101020204000C00000201010101020206002800000201010001020101003C00320201020101020105C35000000301020101020106003200000301020001020104001000320301020101020105C350000003010201010202030028000003010200010201030032003203010201010202059C4000000301020101010206002800000301020001010201003C00320301020101010206003200000301020101010204000C000003010200010101010050003203010201010101059C40000003010201010101030032000003010200010101040010003203010001010101060032000003010001010102030028000003010001010101010050003203010000010102059C4000000301000001010206002800000301000001010010")
	} else {
		data = []byte{0x00, 0x00, 0x00, 0x00}
		s.logger.Info("GET_TENROUIRAI request for unknown type")
	}
	if err != nil {
		panic(err)
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfPostTenrouirai(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfPostTenrouirai)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfGetDailyMissionMaster(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetDailyMissionPersonal(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetDailyMissionPersonal(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEquipSkinHist(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetEquipSkinHist)
	// Transmog / reskin system,  bitmask of 3200 bytes length
	// presumably divided by 5 sections for 5120 armour IDs covered
	// +10,000 for actual ID to be unlocked by each bit
	// Returning 3200 bytes of FF just unlocks everything for now
	var data []byte
	err := s.Server.db.QueryRow("SELECT COALESCE(skin_hist::bytea, $2::bytea) FROM characters WHERE id = $1", s.CharID, make([]byte, 0xC80)).Scan(&data)
	if err != nil {
		s.logger.Fatal("Failed to get skin_hist savedata from db", zap.Error(err))
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfUpdateEquipSkinHist(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateEquipSkinHist)
	// sends a raw armour ID back that needs to be mapped into the persistent bitmask above (-10,000)
	var data []byte
	err := s.Server.db.QueryRow("SELECT COALESCE(skin_hist, $2) FROM characters WHERE id = $1", s.CharID, make([]byte, 0xC80)).Scan(&data)
	if err != nil {
		s.logger.Fatal("Failed to get skin_hist from db", zap.Error(err))
	}

	var bit int
	var startByte int
	switch pkt.MogType {
	case 0: // legs
		bit = int(pkt.ArmourID) - 10000
		startByte = 0
	case 1:
		bit = int(pkt.ArmourID) - 10000
		startByte = 640
	case 2:
		bit = int(pkt.ArmourID) - 10000
		startByte = 1280
	case 3:
		bit = int(pkt.ArmourID) - 10000
		startByte = 1920
	case 4:
		bit = int(pkt.ArmourID) - 10000
		startByte = 2560
	}
	// psql set_bit could also work but I couldn't get it working
	byteInd := (bit / 8)
	bitInByte := bit % 8
	data[startByte+byteInd] |= bits.Reverse8((1 << uint(bitInByte)))
	_, err = s.Server.db.Exec("UPDATE characters SET skin_hist=$1 WHERE id=$2", data, s.CharID)
	if err != nil {
		panic(err)
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfGetUdShopCoin(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdShopCoin)
	data, _ := hex.DecodeString("0000000000000001")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfUseUdShopCoin(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetEnhancedMinidata(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetEnhancedMinidata)
	// this looks to be the detailed chunk of information you can pull up on players in town
	var data []byte
	err := s.Server.db.QueryRow("SELECT minidata FROM characters WHERE id = $1", pkt.CharID).Scan(&data)
	if err != nil {
		data = make([]byte, 0x400) // returning empty might avoid a client softlock
		//s.logger.Fatal("Failed to get minidata from db", zap.Error(err))
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfSetEnhancedMinidata(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSetEnhancedMinidata)
	_, err := s.Server.db.Exec("UPDATE characters SET minidata=$1 WHERE id=$2", pkt.RawDataPayload, s.CharID)
	if err != nil {
		s.logger.Fatal("Failed to update minidata in db", zap.Error(err))
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfGetLobbyCrowd(s *Session, p mhfpacket.MHFPacket) {
	// this requests a specific server's population but seems to have been
	// broken at some point on live as every example response across multiple
	// servers sends back the exact same information?
	// It can be worried about later if we ever get to the point where there are
	// full servers to actually need to migrate people from and empty ones to
	pkt := p.(*mhfpacket.MsgMhfGetLobbyCrowd)
	blankData := make([]byte, 0x320)
	doAckBufSucceed(s, pkt.AckHandle, blankData)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfGetTrendWeapon(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetTrendWeapon)
	// TODO (Fist): Work out actual format limitations, seems to be final upgrade
	// for weapons and it traverses its upgrade tree to recommend base as final
	// 423C correlates with most popular magnet spike in use on JP
	// 2A 00 3C 44 00 3C 76 00 3F EA 01 0F 20 01 0F 50 01 0F F8 02 3C 7E 02 3D
	// F3 02 40 2A 03 3D 65 03 3F 2A 03 40 36 04 3D 59 04 41 E7 04 43 3E 05 0A
	// ED 05 0F 4C 05 0F F2 06 3A FE 06 41 E8 06 41 FA 07 3B 02 07 3F ED 07 40
	// 24 08 3D 37 08 3F 66 08 41 EC 09 3D 38 09 3F 8A 09 41 EE 0A 0E 78 0A 0F
	// AA 0A 0F F9 0B 3E 2E 0B 41 EF 0B 42 FB 0C 41 F0 0C 43 3F 0C 43 EE 0D 41 F1 0D 42 10 0D 42 3C 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	doAckBufSucceed(s, pkt.AckHandle, make([]byte, 0xA9))
}

func handleMsgMhfUpdateUseTrendWeaponLog(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateUseTrendWeaponLog)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}
