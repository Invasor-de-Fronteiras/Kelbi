package channelserver

import (
	"erupe-ce/common/byteframe"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"erupe-ce/network/mhfpacket"
)

func removeSessionFromSemaphore(s *Session) {
	s.Server.semaphoreLock.Lock()
	for _, semaphore := range s.Server.Semaphore {
		if _, exists := semaphore.ReservedClientSlots[s.CharID]; exists {
			delete(semaphore.ReservedClientSlots, s.CharID)
		}
		if _, exists := semaphore.clients[s]; exists {
			delete(semaphore.clients, s)
		}
	}

	s.Server.semaphoreLock.Unlock()
}

func handleMsgSysCreateSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCreateSemaphore)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x03, 0x00, 0x0d})
}

func destructEmptySemaphores(s *Session) {
	s.Server.semaphoreLock.Lock()
	for id, sema := range s.Server.Semaphore {
		if len(sema.ReservedClientSlots) == 0 && len(sema.clients) == 0 {
			s.Server.semaphoreLock.Unlock()
			delete(s.Server.Semaphore, id)
			s.Server.semaphoreLock.Lock()
			if strings.HasPrefix(id, "hs_l0u3B5") {
				releaseRaviSemaphore(s, sema)
			}
			s.logger.Debug("Destructed semaphore", zap.String("sema.id_semaphore", id))
		}
	}
	s.Server.semaphoreLock.Unlock()
}

func releaseRaviSemaphore(s *Session, sema *Semaphore) {
	if !strings.HasSuffix(sema.StageId, "5") {
		delete(sema.ReservedClientSlots, s.CharID)
		delete(sema.clients, s)
	}
	if len(sema.ReservedClientSlots) == 0 && len(sema.clients) == 0 {
		s.logger.Debug("Raviente semaphore is empty, resetting")
		resetRavi(s)
	}
}

func handleMsgSysDeleteSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysDeleteSemaphore)
	sem := pkt.AckHandle
	if s.Server.Semaphore != nil {
		destructEmptySemaphores(s)
		s.Server.semaphoreLock.Lock()
		for id, sema := range s.Server.Semaphore {
			if sema.Id == sem {
				if strings.HasPrefix(id, "hs_l0u3B5") {
					releaseRaviSemaphore(s, sema)
					s.Server.semaphoreLock.Unlock()
					return
				}
				s.Server.semaphoreLock.Unlock()
				delete(s.Server.Semaphore, id)
				s.logger.Debug("Destructed semaphore", zap.String("sema.id_semaphore", id))
				return
			}
		}
		s.Server.semaphoreLock.Unlock()
	}
}

func handleMsgSysCreateAcquireSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCreateAcquireSemaphore)
	SemaphoreID := pkt.SemaphoreID
	newSemaphore, exists := s.Server.Semaphore[SemaphoreID]

	fmt.Printf("Got reserve stage req, StageID: %v\n\n", SemaphoreID)
	if !exists {
		s.Server.semaphoreLock.Lock()
		if strings.HasPrefix(SemaphoreID, "hs_l0u3B5") {
			suffix, _ := strconv.ParseUint(pkt.SemaphoreID[len(pkt.SemaphoreID)-1:], 10, 32)
			s.Server.Semaphore[SemaphoreID] = &Semaphore{
				StageId:             pkt.SemaphoreID,
				Id:                  uint32(suffix),
				clients:             make(map[*Session]uint32),
				ReservedClientSlots: make(map[uint32]interface{}),
				MaxPlayers:          32,
			}
		} else {
			s.Server.Semaphore[SemaphoreID] = NewSemaphore(s.Server, SemaphoreID, 1)
		}
		newSemaphore = s.Server.Semaphore[SemaphoreID]
		s.Server.semaphoreLock.Unlock()
	}

	newSemaphore.Lock()
	defer newSemaphore.Unlock()
	if _, exists := newSemaphore.ReservedClientSlots[s.CharID]; exists {
		bf := byteframe.NewByteFrame()
		bf.WriteUint32(newSemaphore.Id)
		doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
	} else if uint16(len(newSemaphore.ReservedClientSlots)) < newSemaphore.MaxPlayers {
		newSemaphore.ReservedClientSlots[s.CharID] = nil
		newSemaphore.clients[s] = s.CharID
		s.Lock()
		s.Semaphore = newSemaphore
		s.Unlock()
		bf := byteframe.NewByteFrame()
		bf.WriteUint32(newSemaphore.Id)
		doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
	} else {
		doAckSimpleFail(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
	}
}

func handleMsgSysAcquireSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysAcquireSemaphore)
	if sema, exists := s.Server.Semaphore[pkt.SemaphoreID]; exists {
		sema.clients[s] = s.CharID
		bf := byteframe.NewByteFrame()
		bf.WriteUint32(sema.Id)
		doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
	} else {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
	}

}

func handleMsgSysReleaseSemaphore(s *Session, p mhfpacket.MHFPacket) {
	//pkt := p.(*mhfpacket.MsgSysReleaseSemaphore)
}

func handleMsgSysCheckSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCheckSemaphore)
	resp := []byte{0x00, 0x00, 0x00, 0x00}
	s.Server.semaphoreLock.Lock()
	if _, exists := s.Server.Semaphore[pkt.SemaphoreID]; exists {
		resp = []byte{0x00, 0x00, 0x00, 0x01}
	}
	s.Server.semaphoreLock.Unlock()
	doAckSimpleSucceed(s, pkt.AckHandle, resp)
}
