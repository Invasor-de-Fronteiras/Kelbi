package channelserver

import (
	"sync"

	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
)

// Object holds infomation about a specific object.
type Object struct {
	sync.RWMutex
	Id          uint32  `json:"id"`
	OwnerCharID uint32  `json:"ownerCharID"`
	X           float32 `json:"x"`
	Y           float32 `json:"y"`
	Z           float32 `json:"z"`
}

// StageBinaryKey is a struct used as a map key for identifying a stage binary part.
type StageBinaryKey struct {
	Id0 uint8 `json:"id0"`
	Id1 uint8 `json:"id1"`
}

// Stage holds stage-specific information
type Stage struct {
	sync.RWMutex

	// Stage ID string
	Id string `json:"id"`

	// Objects
	Objects     map[uint32]*Object `json:"objects"`
	ObjectIndex uint8              `json:"objectIndex"`

	// Map of session -> charID.
	// These are Clients that are CURRENTLY in the stage
	Clients map[*Session]uint32 `json:"-"`

	// Map of charID -> bool, key represents whether they are ready
	// These are clients that aren't in the stage, but have reserved a slot (for quests, etc).
	ReservedClientSlots map[uint32]bool `json:"reservedClientSlots"`

	// These are raw binary blobs that the stage owner sets,
	// other clients expect the server to echo them back in the exact same format.
	RawBinaryData map[StageBinaryKey][]byte `json:"-"`

	MaxPlayers uint16 `json:"maxPlayers"`
	Password   string `json:"password"`
	CreatedAt  string `json:"createdAt"`
}

// NewStage creates a new stage with intialized values.
func NewStage(ID string) *Stage {
	s := &Stage{
		Id:                  ID,
		Clients:             make(map[*Session]uint32),
		ReservedClientSlots: make(map[uint32]bool),
		Objects:             make(map[uint32]*Object),
		ObjectIndex:         0,
		RawBinaryData:       make(map[StageBinaryKey][]byte),
		MaxPlayers:          4,
		CreatedAt:           time.Now().Format("01-02-2006 15:04:05"),
	}
	return s
}

// BroadcastMHF queues a MHFPacket to be sent to all sessions in the stage.
func (s *Stage) BroadcastMHF(pkt mhfpacket.MHFPacket, ignoredSession *Session) {
	// Broadcast the data.
	for session := range s.Clients {
		if session == ignoredSession {
			continue
		}

		// Make the header
		bf := byteframe.NewByteFrame()
		bf.WriteUint16(uint16(pkt.Opcode()))

		// Build the packet onto the byteframe.
		pkt.Build(bf, session.clientContext)

		// Enqueue in a non-blocking way that drops the packet if the connections send buffer channel is full.
		session.QueueSendNonBlocking(bf.Data())
	}
}

func (s *Stage) isCharInQuestByID(charID uint32) bool {
	if _, exists := s.ReservedClientSlots[charID]; exists {
		return exists
	}

	return false
}

func (s *Stage) isQuest() bool {
	return len(s.ReservedClientSlots) > 0
}

func (s *Stage) GetName() string {
	switch s.Id {
	case MezeportaStageId:
		return "Mezeporta"
	case GuildHallLv1StageId:
		return "Guild Hall Lv1"
	case GuildHallLv2StageId:
		return "Guild Hall Lv2"
	case GuildHallLv3StageId:
		return "Guild Hall Lv3"
	case PugiFarmStageId:
		return "Pugi Farm"
	case RastaBarStageId:
		return "Rasta Bar"
	case PalloneCaravanStageId:
		return "Pallone Caravan"
	case GookFarmStageId:
		return "Gook Farm"
	case DivaFountainStageId:
		return "Diva Fountain"
	case DivaHallStageId:
		return "Diva Hall"
	case MezFesStageId:
		return "Mez Fes"
	default:
		return ""
	}
}

func (s *Stage) NextObjectID() uint32 {
	s.ObjectIndex = s.ObjectIndex + 1
	// Objects beyond 127 do not duplicate correctly
	// Indexes 0 and 127 does not update position correctly
	if s.ObjectIndex == 127 {
		s.ObjectIndex = 1
	}
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(0)
	bf.WriteUint8(s.ObjectIndex)
	bf.WriteUint16(0)
	obj := uint32(bf.Data()[3]) | uint32(bf.Data()[2])<<8 | uint32(bf.Data()[1])<<16 | uint32(bf.Data()[0])<<24
	return obj
}
