package channelserver

import (
	"fmt"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
)

func handleMsgSysCreateObject(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCreateObject)

	// Prevent reusing an object index
	var nextID uint32
	for {
		exists := false
		nextID = s.stage.NextObjectID()
		for _, object := range s.stage.Objects {
			if object.Id == nextID {
				exists = true
				break
			}
		}
		if exists == false {
			break
		}
	}

	s.stage.Lock()
	newObj := &Object{
		Id:          nextID,
		OwnerCharID: s.charID,
		X:           pkt.X,
		Y:           pkt.Y,
		Z:           pkt.Z,
	}
	s.stage.Objects[s.charID] = newObj
	s.stage.Unlock()

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
	s.stage.BroadcastMHF(dupObjUpdate, s)
}

func handleMsgSysDeleteObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysPositionObject(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysPositionObject)

	if s.server.erupeConfig.DevMode && s.server.erupeConfig.DevModeOptions.LogInboundMessages {
		fmt.Printf("[%s - %s] with objectID [%d] move to (%f,%f,%f)\n\n", s.Name, s.stageID, pkt.ObjID, pkt.X, pkt.Y, pkt.Z)
	}
	s.stage.Lock()
	object, ok := s.stage.Objects[s.charID]
	if ok {
		object.X = pkt.X
		object.Y = pkt.Y
		object.Z = pkt.Z
	}
	s.stage.Unlock()
	// One of the few packets we can just re-broadcast directly.
	s.stage.BroadcastMHF(pkt, s)
}

func handleMsgSysRotateObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysDuplicateObject(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysSetObjectBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysSetObjectBinary)
	for _, session := range s.server.Sessions {
		if session.charID == s.charID {
			s.server.userBinaryPartsLock.Lock()
			s.server.userBinaryParts[userBinaryPartID{charID: s.charID, index: 3}] = pkt.RawDataPayload
			s.server.userBinaryPartsLock.Unlock()
			msg := &mhfpacket.MsgSysNotifyUserBinary{
				CharID:     s.charID,
				BinaryType: 3,
			}
			s.server.BroadcastMHF(msg, s)
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
