package channelserver

import (
	"fmt"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
)

func handleMsgSysCreateObject(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCreateObject)

	if s.Stage == nil {
		message := fmt.Sprintf("Player %s failed to create object in Stage  %s", s.Name, s.StageID)
		s.Server.logger.Error(message)
		s.Server.discordBot.LogsChannelSend(message)
		LogoutPlayer(s)
		return
	}

	s.Stage.Lock()
	// Prevent reusing an object index
	var nextID uint32
	for {
		exists := false
		nextID = s.Stage.NextObjectID()
		for _, object := range s.Stage.Objects {
			if object.Id == nextID {
				exists = true
				break
			}
		}
		if !exists {
			break
		}
	}

	newObj := &Object{
		Id:          nextID,
		OwnerCharID: s.CharID,
		X:           pkt.X,
		Y:           pkt.Y,
		Z:           pkt.Z,
	}
	s.Stage.Objects[s.CharID] = newObj
	s.Stage.Unlock()

	// Response to our requesting client.
	resp := byteframe.NewByteFrame()
	resp.WriteUint32(newObj.Id) // New local obj handle.
	doAckSimpleSucceed(s, pkt.AckHandle, resp.Data())
	// Duplicate the object creation to all sessions in the same stage.
	dupObjUpdate := &mhfpacket.MsgSysDuplicateObject{
		ObjID:       newObj.Id,
		X:           newObj.X,
		Y:           newObj.Y,
		Z:           newObj.Z,
		OwnerCharID: newObj.OwnerCharID,
	}

	s.logger.Info(fmt.Sprintf("Broadcasting new object: %s (%d)", s.Name, newObj.Id))
	s.Stage.BroadcastMHF(dupObjUpdate, s)
}

func handleMsgSysDeleteObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysPositionObject(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysPositionObject)
	if s.Server.erupeConfig.DevMode && s.Server.erupeConfig.DevModeOptions.LogInboundMessages {
		fmt.Printf("[%s - %s] with objectID [%d] move to (%f,%f,%f)\n\n", s.Name, s.StageID, pkt.ObjID, pkt.X, pkt.Y, pkt.Z)
	}
	s.Stage.Lock()
	object, ok := s.Stage.Objects[s.CharID]
	if ok {
		object.X = pkt.X
		object.Y = pkt.Y
		object.Z = pkt.Z
	}
	s.Stage.Unlock()
	// One of the few packets we can just re-broadcast directly.
	s.Stage.BroadcastMHF(pkt, s)
}

func handleMsgSysRotateObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDuplicateObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetObjectBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetObjectBinary)
	for _, session := range s.Server.Sessions {
		if session.CharID == s.CharID {
			s.Server.userBinaryPartsLock.Lock()
			s.Server.userBinaryParts[userBinaryPartID{charID: s.CharID, index: 3}] = pkt.RawDataPayload
			s.Server.userBinaryPartsLock.Unlock()
			msg := &mhfpacket.MsgSysNotifyUserBinary{
				CharID:     s.CharID,
				BinaryType: 3,
			}
			s.Server.BroadcastMHF(msg, s)
		}
	}
}

func handleMsgSysGetObjectBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysGetObjectOwner(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysUpdateObjectBinary(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysCleanupObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysAddObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDelObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDispObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysHideObject(s *Session, p mhfpacket.MHFPacket) {}
