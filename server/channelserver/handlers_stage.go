package channelserver

import (
	"fmt"
	"strings"
	"time"

	"erupe-ce/common/byteframe"
	ps "erupe-ce/common/pascalstring"
	"erupe-ce/network/mhfpacket"

	"go.uber.org/zap"
)

func handleMsgSysCreateStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCreateStage)
	s.Server.Lock()
	defer s.Server.Unlock()
	if _, exists := s.Server.Stages[pkt.StageID]; exists {
		doAckSimpleFail(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
	} else {
		stage := NewStage(pkt.StageID)
		stage.MaxPlayers = uint16(pkt.PlayerCount)
		s.Server.Stages[stage.Id] = stage
		doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
	}
}

func handleMsgSysStageDestruct(s *Session, p mhfpacket.MHFPacket) {}

func doStageTransfer(s *Session, ackHandle uint32, stageID string) {
	s.Server.Lock()
	stage, exists := s.Server.Stages[stageID]
	s.Server.Unlock()

	if exists {
		stage.Lock()
		stage.Clients[s] = s.CharID
		stage.Unlock()
	} else { // Create new stage object
		s.Server.Lock()
		s.Server.Stages[stageID] = NewStage(stageID)
		stage = s.Server.Stages[stageID]
		s.Server.Unlock()
		stage.Lock()
		stage.Clients[s] = s.CharID
		stage.Unlock()
	}

	// Ensure this session no longer belongs to reservations.
	if s.Stage != nil {
		removeSessionFromStage(s)
	}

	// Save our new stage ID and pointer to the new stage itself.
	s.Lock()
	s.StageID = stageID
	s.Stage = s.Server.Stages[stageID]
	s.Unlock()

	// Tell the client to cleanup its current stage objects.
	s.QueueSendMHF(&mhfpacket.MsgSysCleanupObject{})

	// Confirm the stage entry.
	doAckSimpleSucceed(s, ackHandle, []byte{0x00, 0x00, 0x00, 0x00})

	var temp mhfpacket.MHFPacket
	newNotif := byteframe.NewByteFrame()

	// Cast existing user data to new user
	if !s.UserEnteredStage {
		s.UserEnteredStage = true

		for _, session := range s.Server.Sessions {
			if s == session {
				continue
			}
			temp = &mhfpacket.MsgSysInsertUser{CharID: session.CharID}
			newNotif.WriteUint16(uint16(temp.Opcode()))
			// nolint:errcheck
			temp.Build(newNotif, s.clientContext)
			for i := 0; i < 3; i++ {
				temp = &mhfpacket.MsgSysNotifyUserBinary{
					CharID:     session.CharID,
					BinaryType: uint8(i + 1),
				}
				newNotif.WriteUint16(uint16(temp.Opcode()))
				// nolint:errcheck
				temp.Build(newNotif, s.clientContext)
			}
		}
	}

	if s.Stage != nil { // avoids lock up when using bed for dream quests
		// Notify the client to duplicate the existing objects.
		s.logger.Info(fmt.Sprintf("Sending existing stage objects to %s", s.Name))
		s.Stage.RLock()
		var temp mhfpacket.MHFPacket
		for _, obj := range s.Stage.Objects {
			if obj.OwnerCharID == s.CharID {
				continue
			}
			temp = &mhfpacket.MsgSysDuplicateObject{
				ObjID:       obj.Id,
				X:           obj.X,
				Y:           obj.Y,
				Z:           obj.Z,
				Unk0:        0,
				OwnerCharID: obj.OwnerCharID,
			}
			newNotif.WriteUint16(uint16(temp.Opcode()))
			// nolint:errcheck
			temp.Build(newNotif, s.clientContext)
		}
		s.Stage.RUnlock()
	}

	newNotif.WriteUint16(0x0010) // End it.
	if len(newNotif.Data()) > 2 {
		s.QueueSend(newNotif.Data())
	}
}

func destructEmptyStages(s *Session) {
	s.Server.Lock()
	defer s.Server.Unlock()
	for _, stage := range s.Server.Stages {
		// Destroy empty Quest/My series/Guild stages.
		if stage.Id[3:5] == "Qs" || stage.Id[3:5] == "Ms" || stage.Id[3:5] == "Gs" || stage.Id[3:5] == "Ls" {
			if len(stage.ReservedClientSlots) == 0 && len(stage.Clients) == 0 {
				delete(s.Server.Stages, stage.Id)
				s.logger.Debug("Destructed stage", zap.String("stage.id", stage.Id))
			}
		}
	}
}

func removeSessionFromStage(s *Session) {
	// Remove client from old stage.
	s.Stage.Lock()
	delete(s.Stage.Clients, s)

	// Delete old stage objects owned by the client.
	s.logger.Info("Sending notification to old stage clients")
	for _, object := range s.Stage.Objects {
		if object.OwnerCharID == s.CharID {
			s.Stage.BroadcastMHF(&mhfpacket.MsgSysDeleteObject{ObjID: object.Id}, s)
			delete(s.Stage.Objects, object.OwnerCharID)
		}
	}
	s.Stage.Unlock()
	destructEmptyStages(s)
	destructEmptySemaphores(s)
}

func handleMsgSysEnterStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysEnterStage)

	// Push our current stage ID to the movement stack before entering another one.
	if s.StageID == "" {
		s.stageMoveStack.Set(pkt.StageID)
	} else {
		// s.Stage.ReservedClientSlots[s.CharID] = false
		s.stageMoveStack.Push(s.StageID)
		s.stageMoveStack.Lock()
	}

	s.QueueSendMHF(&mhfpacket.MsgSysCleanupObject{})
	if s.ReservationStage != nil {
		s.ReservationStage = nil
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

	delete(s.Stage.ReservedClientSlots, s.CharID)
	delete(s.Server.Stages[backStage].ReservedClientSlots, s.CharID)

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
	// I think this is supposed to mark a stage as no longer able to accept client reservations
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgSysUnlockStage(s *Session, p mhfpacket.MHFPacket) {
	if s.ReservationStage != nil {
		s.ReservationStage.RLock()
		defer s.ReservationStage.RUnlock()

		for charID := range s.ReservationStage.ReservedClientSlots {
			session := s.Server.FindSessionByCharID(charID)
			session.QueueSendMHF(&mhfpacket.MsgSysStageDestruct{})
		}

		delete(s.Server.Stages, s.ReservationStage.Id)
	}

	destructEmptyStages(s)
}

func handleMsgSysReserveStage(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysReserveStage)
	if stage, exists := s.Server.Stages[pkt.StageID]; exists {
		stage.Lock()
		defer stage.Unlock()
		if _, exists := stage.ReservedClientSlots[s.CharID]; exists {
			switch pkt.Ready {
			case 1: // 0x01
				stage.ReservedClientSlots[s.CharID] = false
			case 17: // 0x11
				stage.ReservedClientSlots[s.CharID] = true
			}
			doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
		} else if uint16(len(stage.ReservedClientSlots)) < stage.MaxPlayers {
			if len(stage.Password) > 0 {
				if stage.Password != s.StagePass {
					doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
					return
				}
			}
			stage.ReservedClientSlots[s.CharID] = false
			// Save the reservation stage in the session for later use in MsgSysUnreserveStage.
			s.Lock()
			s.ReservationStage = stage
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
	stage := s.ReservationStage
	s.ReservationStage = nil
	s.Unlock()
	if stage != nil {
		stage.Lock()
		// nolint:gosimple
		if _, exists := stage.ReservedClientSlots[s.CharID]; exists {
			delete(stage.ReservedClientSlots, s.CharID)
		}
		stage.Unlock()
	}
}

func handleMsgSysSetStagePass(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetStagePass)
	s.Lock()
	stage := s.ReservationStage
	s.Unlock()
	if stage != nil {
		stage.Lock()
		// Will only exist if host.
		if _, exists := stage.ReservedClientSlots[s.CharID]; exists {
			stage.Password = pkt.Password
		}
		stage.Unlock()
	} else {
		// Store for use on next ReserveStage.
		s.Lock()
		s.StagePass = pkt.Password
		s.Unlock()
	}
}

func handleMsgSysSetStageBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetStageBinary)
	if stage, exists := s.Server.Stages[pkt.StageID]; exists {
		stage.Lock()
		stage.RawBinaryData[StageBinaryKey{pkt.BinaryType0, pkt.BinaryType1}] = pkt.RawDataPayload
		stage.Unlock()
	} else {
		s.logger.Warn("Failed to get stage", zap.String("StageID", pkt.StageID))
	}
}

func handleMsgSysGetStageBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysGetStageBinary)
	if stage, exists := s.Server.Stages[pkt.StageID]; exists {
		stage.Lock()
		if binaryData, exists := stage.RawBinaryData[StageBinaryKey{pkt.BinaryType0, pkt.BinaryType1}]; exists {
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
	if stage, exists := s.Server.Stages[pkt.StageID]; exists {
		if pkt.BinaryType0 == 1 && pkt.BinaryType1 == 12 {
			// This might contain the hunter count, or max player count?
			doAckBufSucceed(s, pkt.AckHandle, []byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			return
		}
		for {
			s.logger.Debug("MsgSysWaitStageBinary before lock and get stage")
			stage.Lock()
			stageBinary, gotBinary := stage.RawBinaryData[StageBinaryKey{pkt.BinaryType0, pkt.BinaryType1}]
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
	s.Server.stagesLock.RLock()
	defer s.Server.stagesLock.RUnlock()

	// Build the response
	resp := byteframe.NewByteFrame()
	bf := byteframe.NewByteFrame()
	var joinable int
	for sid, stage := range s.Server.Stages {
		stage.RLock()
		defer stage.RUnlock()

		if len(stage.ReservedClientSlots) == 0 && len(stage.Clients) == 0 {
			continue
		}

		if !strings.Contains(stage.Id, pkt.StagePrefix) {
			continue
		}

		joinable++

		resp.WriteUint16(uint16(len(stage.ReservedClientSlots))) // Reserved players.
		resp.WriteUint16(0)                                      // Unk
		resp.WriteUint8(0)                                       // Unk
		resp.WriteBool(len(stage.Clients) > 0)                   // Has departed.
		resp.WriteUint16(stage.MaxPlayers)                       // Max players.
		if len(stage.Password) > 0 {
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
