package channelserver

import (
	"fmt"
	"time"

	"erupe-ce/common/byteframe"
	ps "erupe-ce/common/pascalstring"
	"erupe-ce/network/mhfpacket"

	"go.uber.org/zap"
)

func handleMsgSysCreateStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCreateStage)
	s.server.Lock()
	defer s.server.Unlock()
	if _, exists := s.server.stages[pkt.StageID]; exists {
		doAckSimpleFail(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
	} else {
		stage := NewStage(pkt.StageID)
		stage.maxPlayers = uint16(pkt.PlayerCount)
		s.server.stages[stage.id] = stage
		doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
	}
}

func handleMsgSysStageDestruct(s *Session, p mhfpacket.MHFPacket) {}

func doStageTransfer(s *Session, ackHandle uint32, stageID string) {
	s.server.Lock()
	stage, exists := s.server.stages[stageID]
	s.server.Unlock()

	if exists {
		stage.Lock()
		stage.clients[s] = s.charID
		stage.Unlock()
	} else { // Create new stage object
		s.server.Lock()
		s.server.stages[stageID] = NewStage(stageID)
		stage = s.server.stages[stageID]
		s.server.Unlock()
		stage.Lock()
		stage.clients[s] = s.charID
		stage.Unlock()
	}

	// Ensure this session no longer belongs to reservations.
	if s.stage != nil {
		removeSessionFromStage(s)
	}

	// Save our new stage ID and pointer to the new stage itself.
	s.Lock()
	s.stageID = stageID
	s.stage = s.server.stages[stageID]
	s.Unlock()

	// Tell the client to cleanup its current stage objects.
	s.QueueSendMHF(&mhfpacket.MsgSysCleanupObject{})

	// Confirm the stage entry.
	doAckSimpleSucceed(s, ackHandle, []byte{0x00, 0x00, 0x00, 0x00})

	var temp mhfpacket.MHFPacket
	newNotif := byteframe.NewByteFrame()

	// Cast existing user data to new user
	if !s.userEnteredStage {
		s.userEnteredStage = true

		for _, session := range s.server.sessions {
			if s == session {
				continue
			}
			temp = &mhfpacket.MsgSysInsertUser{CharID: session.charID}
			newNotif.WriteUint16(uint16(temp.Opcode()))
			temp.Build(newNotif, s.clientContext)
			for i := 0; i < 3; i++ {
				temp = &mhfpacket.MsgSysNotifyUserBinary{
					CharID:     session.charID,
					BinaryType: uint8(i + 1),
				}
				newNotif.WriteUint16(uint16(temp.Opcode()))
				temp.Build(newNotif, s.clientContext)
			}
		}
	}

	if s.stage != nil { // avoids lock up when using bed for dream quests
		// Notify the client to duplicate the existing objects.
		s.logger.Info(fmt.Sprintf("Sending existing stage objects to %s", s.Name))
		s.stage.RLock()
		var temp mhfpacket.MHFPacket
		for _, obj := range s.stage.objects {
			if obj.ownerCharID == s.charID {
				continue
			}
			temp = &mhfpacket.MsgSysDuplicateObject{
				ObjID:       obj.id,
				X:           obj.x,
				Y:           obj.y,
				Z:           obj.z,
				Unk0:        0,
				OwnerCharID: obj.ownerCharID,
			}
			newNotif.WriteUint16(uint16(temp.Opcode()))
			temp.Build(newNotif, s.clientContext)
		}
		s.stage.RUnlock()
	}

	newNotif.WriteUint16(0x0010) // End it.
	if len(newNotif.Data()) > 2 {
		s.QueueSend(newNotif.Data())
	}
}

func destructEmptyStages(s *Session) {
	s.server.Lock()
	defer s.server.Unlock()
	for _, stage := range s.server.stages {
		// Destroy empty Quest/My series/Guild stages.
		if stage.id[3:5] == "Qs" || stage.id[3:5] == "Ms" || stage.id[3:5] == "Gs" {
			if len(stage.reservedClientSlots) == 0 && len(stage.clients) == 0 {
				delete(s.server.stages, stage.id)
				s.logger.Debug("Destructed stage", zap.String("stage.id", stage.id))
			}
		}
	}
}

func removeSessionFromStage(s *Session) {
	// Remove client from old stage.
	s.stage.Lock()
	delete(s.stage.clients, s)
	delete(s.stage.reservedClientSlots, s.charID)

	// Delete old stage objects owned by the client.
	s.logger.Info("Sending notification to old stage clients")
	for _, object := range s.stage.objects {
		if object.ownerCharID == s.charID {
			s.stage.BroadcastMHF(&mhfpacket.MsgSysDeleteObject{ObjID: object.id}, s)
			delete(s.stage.objects, object.ownerCharID)
		}
	}
	s.stage.Unlock()
	destructEmptyStages(s)
	destructEmptySemaphores(s)
}

func handleMsgSysEnterStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysEnterStage)

	// Push our current stage ID to the movement stack before entering another one.
	if s.stageID == "" {
		s.stageMoveStack.Set(pkt.StageID)
	} else {
		s.stageMoveStack.Push(s.stageID)
		s.stageMoveStack.Lock()
	}

	s.QueueSendMHF(&mhfpacket.MsgSysCleanupObject{})
	if s.reservationStage != nil {
		s.reservationStage = nil
	}

	doStageTransfer(s, pkt.AckHandle, pkt.StageID)
}

func handleMsgSysBackStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysBackStage)

	// Transfer back to the saved stage ID before the previous move or enter.
	s.stageMoveStack.Unlock()
	backStage, err := s.stageMoveStack.Pop()

	if err != nil {
		panic(err)
	}

	doStageTransfer(s, pkt.AckHandle, backStage)
}

func handleMsgSysMoveStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysMoveStage)

	// Set a new move stack from the given stage ID if unlocked
	if !s.stageMoveStack.Locked {
		s.stageMoveStack.Set(pkt.StageID)
	}

	doStageTransfer(s, pkt.AckHandle, pkt.StageID)
}

func handleMsgSysLeaveStage(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysLockStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysLockStage)
	// TODO(Andoryuuta): What does this packet _actually_ do?
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgSysUnlockStage(s *Session, p mhfpacket.MHFPacket) {
	s.reservationStage.RLock()
	defer s.reservationStage.RUnlock()

	for charID := range s.reservationStage.reservedClientSlots {
		session := s.server.FindSessionByCharID(charID)
		session.QueueSendMHF(&mhfpacket.MsgSysStageDestruct{})
	}

	s.server.Lock()
	defer s.server.Unlock()

	delete(s.server.stages, s.reservationStage.id)
}

func handleMsgSysReserveStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysReserveStage)

	if stage, exists := s.server.stages[pkt.StageID]; exists {
		stage.Lock()
		defer stage.Unlock()
		if _, exists := stage.reservedClientSlots[s.charID]; exists {
			switch pkt.Ready {
			case 1: // 0x01
				stage.reservedClientSlots[s.charID] = false
			case 17: // 0x11
				stage.reservedClientSlots[s.charID] = true
			}
			doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
		} else if uint16(len(stage.reservedClientSlots)) < stage.maxPlayers {
			if len(stage.password) > 0 {
				if stage.password != s.stagePass {
					doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
					return
				}
			}
			stage.reservedClientSlots[s.charID] = false
			// Save the reservation stage in the session for later use in MsgSysUnreserveStage.
			s.Lock()
			s.reservationStage = stage
			s.Unlock()
			doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
		} else {
			doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		}
	} else {
		s.logger.Error("Failed to get stage", zap.String("StageID", pkt.StageID))
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
	}
}

func handleMsgSysUnreserveStage(s *Session, p mhfpacket.MHFPacket) {
	s.Lock()
	stage := s.reservationStage
	s.reservationStage = nil
	s.Unlock()
	if stage != nil {
		stage.Lock()
		if _, exists := stage.reservedClientSlots[s.charID]; exists {
			delete(stage.reservedClientSlots, s.charID)
		}
		stage.Unlock()
	}
}

func handleMsgSysSetStagePass(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetStagePass)
	s.Lock()
	stage := s.reservationStage
	s.Unlock()
	if stage != nil {
		stage.Lock()
		// Will only exist if host.
		if _, exists := stage.reservedClientSlots[s.charID]; exists {
			stage.password = pkt.Password
		}
		stage.Unlock()
	} else {
		// Store for use on next ReserveStage.
		s.Lock()
		s.stagePass = pkt.Password
		s.Unlock()
	}
}

func handleMsgSysSetStageBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetStageBinary)
	if stage, exists := s.server.stages[pkt.StageID]; exists {
		stage.Lock()
		stage.rawBinaryData[stageBinaryKey{pkt.BinaryType0, pkt.BinaryType1}] = pkt.RawDataPayload
		stage.Unlock()
	} else {
		s.logger.Warn("Failed to get stage", zap.String("StageID", pkt.StageID))
	}
}

func handleMsgSysGetStageBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysGetStageBinary)

	if stage, exists := s.server.stages[pkt.StageID]; exists {
		stage.Lock()
		if binaryData, exists := stage.rawBinaryData[stageBinaryKey{pkt.BinaryType0, pkt.BinaryType1}]; exists {
			doAckBufSucceed(s, pkt.AckHandle, binaryData)
		} else if pkt.BinaryType1 == 4 {
			// Unknown binary type that is supposedly generated server side
			// Temporary response
			doAckBufSucceed(s, pkt.AckHandle, []byte{})
		} else {
			s.logger.Warn("Failed to get stage binary", zap.Uint8("BinaryType0", pkt.BinaryType0), zap.Uint8("pkt.BinaryType1", pkt.BinaryType1))
			s.logger.Warn("Sending blank stage binary")
			doAckBufSucceed(s, pkt.AckHandle, []byte{})
		}
		stage.Unlock()
	} else {
		s.logger.Warn("Failed to get stage", zap.String("StageID", pkt.StageID))
	}
	s.logger.Debug("MsgSysGetStageBinary Done!")
}

func handleMsgSysWaitStageBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysWaitStageBinary)
	if stage, exists := s.server.stages[pkt.StageID]; exists {
		if pkt.BinaryType0 == 1 && pkt.BinaryType1 == 12 {
			// This might contain the hunter count, or max player count?
			doAckBufSucceed(s, pkt.AckHandle, []byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			return
		}
		for {
			s.logger.Debug("MsgSysWaitStageBinary before lock and get stage")
			stage.Lock()
			stageBinary, gotBinary := stage.rawBinaryData[stageBinaryKey{pkt.BinaryType0, pkt.BinaryType1}]
			stage.Unlock()
			s.logger.Debug("MsgSysWaitStageBinary after lock and get stage")
			if gotBinary {
				doAckBufSucceed(s, pkt.AckHandle, stageBinary)
				break
			} else {
				s.logger.Debug("Waiting stage binary", zap.Uint8("BinaryType0", pkt.BinaryType0), zap.Uint8("pkt.BinaryType1", pkt.BinaryType1))
				time.Sleep(1 * time.Second)
				continue
			}
		}
	} else {
		s.logger.Warn("Failed to get stage", zap.String("StageID", pkt.StageID))
	}
	s.logger.Debug("MsgSysWaitStageBinary Done!")
}

func handleMsgSysEnumerateStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysEnumerateStage)

	// Read-lock the server stage map.
	s.server.stagesLock.RLock()
	defer s.server.stagesLock.RUnlock()

	// Build the response
	resp := byteframe.NewByteFrame()
	bf := byteframe.NewByteFrame()
	var joinable int
	for sid, stage := range s.server.stages {
		stage.RLock()
		defer stage.RUnlock()

		if len(stage.reservedClientSlots) == 0 && len(stage.clients) == 0 {
			continue
		}

		// Check for valid stage type
		if sid[3:5] != "Qs" && sid[3:5] != "Ms" {
			continue
		}

		joinable++

		resp.WriteUint16(uint16(len(stage.reservedClientSlots))) // Reserved players.
		resp.WriteUint16(0)                                      // Unk
		resp.WriteUint8(0)                                       // Unk
		resp.WriteBool(len(stage.clients) > 0)                   // Has departed.
		resp.WriteUint16(stage.maxPlayers)                       // Max players.
		if len(stage.password) > 0 {
			// This byte has also been seen as 1
			// The quest is also recognised as locked when this is 2
			resp.WriteUint8(3)
		} else {
			resp.WriteUint8(0)
		}
		ps.Uint8(resp, sid, false)
	}
	bf.WriteUint16(uint16(joinable))
	bf.WriteBytes(resp.Data())

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}
