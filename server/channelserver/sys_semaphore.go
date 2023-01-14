package channelserver

import (
	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"

	"sync"
)

// Stage holds stage-specific information
type Semaphore struct {
	sync.RWMutex

	// Stage ID string, OLD id_semaphore
	StageId string `json:"stageID"`

	Id uint32 `json:"id"`

	id uint32

	// Map of session -> charID.
	// These are clients that are CURRENTLY in the stage
	clients map[*Session]uint32

	// Map of charID -> interface{}, only the key is used, value is always nil.
	ReservedClientSlots map[uint32]interface{} `json:"reservedClientSlots"`

	// Max Players for Semaphore
	MaxPlayers uint16 `json:"maxPlayers"`
}

// NewStage creates a new stage with intialized values.
func NewSemaphore(s *Server, ID string, MaxPlayers uint16) *Semaphore {
	sema := &Semaphore{
		StageId:             ID,
		Id:                  s.NextSemaphoreID(),
		clients:             make(map[*Session]uint32),
		ReservedClientSlots: make(map[uint32]interface{}),
		MaxPlayers:          MaxPlayers,
	}
	return sema
}

func (s *Semaphore) BroadcastRavi(pkt mhfpacket.MHFPacket) {
	// Broadcast the data.
	for session := range s.clients {

		// Make the header
		bf := byteframe.NewByteFrame()
		bf.WriteUint16(uint16(pkt.Opcode()))

		// Build the packet onto the byteframe.
		// nolint:errcheck // Error return value of `pkt.Build` is not checked (errcheck) pkt.Build(bf, session.clientContext)
		pkt.Build(bf, session.clientContext)

		// Enqueue in a non-blocking way that drops the packet if the connections send buffer channel is full.
		session.QueueSendNonBlocking(bf.Data())
	}
}

// BroadcastMHF queues a MHFPacket to be sent to all sessions in the stage.
func (s *Semaphore) BroadcastMHF(pkt mhfpacket.MHFPacket, ignoredSession *Session) {
	// Broadcast the data.
	for session := range s.clients {
		if session == ignoredSession {
			continue
		}

		// Make the header
		bf := byteframe.NewByteFrame()
		bf.WriteUint16(uint16(pkt.Opcode()))

		// Build the packet onto the byteframe.
		// nolint:errcheck
		pkt.Build(bf, session.clientContext)

		// Enqueue in a non-blocking way that drops the packet if the connections send buffer channel is full.
		session.QueueSendNonBlocking(bf.Data())
	}
}
