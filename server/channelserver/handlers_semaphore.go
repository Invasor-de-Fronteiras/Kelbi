package channelserver

import (
	"erupe-ce/common/byteframe"
	"fmt"
	"go.uber.org/zap"
	"strings"

	"erupe-ce/network/mhfpacket"
)

func removeSessionFromSemaphore(s *Session) {
	s.server.semaphoreLock.Lock()
	for _, semaphore := range s.server.semaphore {
		if _, exists := semaphore.reservedClientSlots[s.charID]; exists {
			delete(semaphore.reservedClientSlots, s.charID)
		}
		if _, exists := semaphore.clients[s]; exists {
			delete(semaphore.clients, s)
		}
	}
	s.server.semaphoreLock.Unlock()
}

func handleMsgSysCreateSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCreateSemaphore)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x03, 0x00, 0x0d})
}

func destructEmptySemaphores(s *Session) {
	s.server.semaphoreLock.Lock()
	for id, sema := range s.server.semaphore {
		if len(sema.reservedClientSlots) == 0 && len(sema.clients) == 0 {
			s.server.semaphoreLock.Unlock()
			delete(s.server.semaphore, id)
			s.server.semaphoreLock.Lock()
			if strings.HasPrefix(id, "hs_l0u3B51") {
				releaseRaviSemaphore(s, sema)
			}
			s.logger.Debug("Destructed semaphore", zap.String("sema.id_semaphore", id))
		}
	}
	s.server.semaphoreLock.Unlock()
}

func releaseRaviSemaphore(s *Session, sema *Semaphore) {
	if !strings.HasSuffix(sema.id_semaphore, "5") {
		delete(sema.reservedClientSlots, s.charID)
		delete(sema.clients, s)
	}
	if len(sema.reservedClientSlots) == 0 && len(sema.clients) == 0 {
		s.logger.Debug("Raviente semaphore is empty, resetting")
		resetRavi(s)
	}
}

func handleMsgSysDeleteSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysDeleteSemaphore)
	sem := pkt.AckHandle
	if s.server.semaphore != nil {
		destructEmptySemaphores(s)
		s.server.semaphoreLock.Lock()
		for id, sema := range s.server.semaphore {
			if sema.id == sem {
				if strings.HasPrefix(id, "hs_l0u3B51") {
					releaseRaviSemaphore(s, sema)
					s.server.semaphoreLock.Unlock()
					return
				}
				s.server.semaphoreLock.Unlock()
				delete(s.server.semaphore, id)
				s.logger.Debug("Destructed semaphore", zap.String("sema.id_semaphore", id))
				return
			}
		}
		s.server.semaphoreLock.Unlock()
	}
}

func handleMsgSysCreateAcquireSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCreateAcquireSemaphore)
	SemaphoreID := pkt.SemaphoreID

	newSemaphore, exists := s.server.semaphore[SemaphoreID]

	fmt.Printf("Got reserve stage req, StageID: %v\n\n", SemaphoreID)
	if !exists {
		s.server.semaphoreLock.Lock()
		if strings.HasPrefix(SemaphoreID, "hs_l0u3B51") {
			s.server.semaphore[SemaphoreID] = NewSemaphore(s.server, SemaphoreID, 32)
			if strings.HasSuffix(SemaphoreID, "3") {
				s.server.raviente.state.semaphoreID = s.server.semaphore[SemaphoreID].id
			} else if strings.HasSuffix(SemaphoreID, "4") {
				s.server.raviente.support.semaphoreID = s.server.semaphore[SemaphoreID].id
			} else if strings.HasSuffix(SemaphoreID, "5") {
				s.server.raviente.register.semaphoreID = s.server.semaphore[SemaphoreID].id
			}
		} else {
			s.server.semaphore[SemaphoreID] = NewSemaphore(s.server, SemaphoreID, 1)
		}
		newSemaphore = s.server.semaphore[SemaphoreID]
		s.server.semaphoreLock.Unlock()
	}

	newSemaphore.Lock()
	defer newSemaphore.Unlock()
	if _, exists := newSemaphore.reservedClientSlots[s.charID]; exists {
		bf := byteframe.NewByteFrame()
		bf.WriteUint32(newSemaphore.id)
		doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
	} else if uint16(len(newSemaphore.reservedClientSlots)) < newSemaphore.maxPlayers {
		newSemaphore.reservedClientSlots[s.charID] = nil
		newSemaphore.clients[s] = s.charID
		s.Lock()
		s.semaphore = newSemaphore
		s.Unlock()
		bf := byteframe.NewByteFrame()
		bf.WriteUint32(newSemaphore.id)
		doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
	} else {
		doAckSimpleFail(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
	}
}

func handleMsgSysAcquireSemaphore(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysReleaseSemaphore(s *Session, p mhfpacket.MHFPacket) {
	//pkt := p.(*mhfpacket.MsgSysReleaseSemaphore)
}

func handleMsgSysCheckSemaphore(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCheckSemaphore)
	resp := []byte{0x00, 0x00, 0x00, 0x00}
	s.server.semaphoreLock.Lock()
	if _, exists := s.server.semaphore[pkt.StageID]; exists {
		resp = []byte{0x00, 0x00, 0x00, 0x01}
	}
	s.server.semaphoreLock.Unlock()
	doAckSimpleSucceed(s, pkt.AckHandle, resp)
}
