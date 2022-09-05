package channelserver

import (
	"fmt"
	"net"
	"sync"

	"erupe-ce/common/byteframe"
	ps "erupe-ce/common/pascalstring"
	"erupe-ce/config"
	"erupe-ce/network/binpacket"
	"erupe-ce/network/mhfpacket"
	"erupe-ce/server/discordbot"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type StageIdType = string

const (
	// GlobalStage is the stage that is used for all users.
	MezeportaStageId      StageIdType = "sl1Ns200p0a0u0"
	GuildHallLv1StageId   StageIdType = "sl1Ns202p0a0u0"
	GuildHallLv2StageId   StageIdType = "sl1Ns203p0a0u0"
	GuildHallLv3StageId   StageIdType = "sl1Ns204p0a0u0"
	PugiFarmStageId       StageIdType = "sl1Ns205p0a0u0"
	RastaBarStageId       StageIdType = "sl1Ns211p0a0u0"
	PalloneCaravanStageId StageIdType = "sl1Ns260p0a0u0"
	GookFarmStageId       StageIdType = "sl1Ns265p0a0u0"
	DivaFountainStageId   StageIdType = "sl2Ns379p0a0u0"
	DivaHallStageId       StageIdType = "sl1Ns445p0a0u0"
	MezFesStageId         StageIdType = "sl1Ns462p0a0u0"
)

// Config struct allows configuring the server.
type Config struct {
	ID          uint16
	Logger      *zap.Logger
	DB          *sqlx.DB
	DiscordBot  *discordbot.DiscordBot
	ErupeConfig *config.Config
	Name        string
	Enable      bool
}

// Map key type for a user binary part.
type userBinaryPartID struct {
	charID uint32
	index  uint8
}

// Server is a MHF channel server.
type Server struct {
	sync.Mutex
	Channels       []*Server `json:"-"`
	ID             uint16    `json:"id"`
	IP             string    `json:"ip"`
	Port           uint16    `json:"port"`
	logger         *zap.Logger
	db             *sqlx.DB
	erupeConfig    *config.Config
	acceptConns    chan net.Conn
	deleteConns    chan net.Conn
	Sessions       map[net.Conn]*Session `json:"-"`
	listener       net.Listener          // Listener that is created when Server.Start is called.
	IsShuttingDown bool                  `json:"isShuttingDown"`

	stagesLock sync.RWMutex
	Stages     map[string]*Stage

	// UserBinary
	userBinaryPartsLock sync.RWMutex
	userBinaryParts     map[userBinaryPartID][]byte

	// Semaphore
	semaphoreLock  sync.RWMutex
	Semaphore      map[string]*Semaphore `json:"semaphore"`
	SemaphoreIndex uint32                `json:"semaphoreIndex"`

	// Discord chat integration
	discordBot *discordbot.DiscordBot

	Name   string `json:"name"`
	Enable bool   `json:"enable"`

	raviente *Raviente
}

type Raviente struct {
	sync.Mutex

	register *RavienteRegister
	state    *RavienteState
	support  *RavienteSupport
}

type RavienteRegister struct {
	nextTime     uint32
	startTime    uint32
	postTime     uint32
	killedTime   uint32
	ravienteType uint32
	maxPlayers   uint32
	carveQuest   uint32
	register     []uint32
}

type RavienteState struct {
	damageMultiplier uint32
	stateData        []uint32
}

type RavienteSupport struct {
	supportData []uint32
}

// Set up the Raviente variables for the server
func NewRaviente() *Raviente {
	ravienteRegister := &RavienteRegister{
		nextTime:     0,
		startTime:    0,
		killedTime:   0,
		postTime:     0,
		ravienteType: 0,
		maxPlayers:   0,
		carveQuest:   0,
	}
	ravienteState := &RavienteState{
		damageMultiplier: 1,
	}
	ravienteSupport := &RavienteSupport{}
	ravienteRegister.register = []uint32{0, 0, 0, 0, 0}
	ravienteState.stateData = []uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ravienteSupport.supportData = []uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	raviente := &Raviente{
		register: ravienteRegister,
		state:    ravienteState,
		support:  ravienteSupport,
	}
	return raviente
}

// NewServer creates a new Server type.
func NewServer(config *Config) *Server {
	s := &Server{
		ID:              config.ID,
		logger:          config.Logger,
		db:              config.DB,
		erupeConfig:     config.ErupeConfig,
		acceptConns:     make(chan net.Conn),
		deleteConns:     make(chan net.Conn),
		Sessions:        make(map[net.Conn]*Session),
		Stages:          make(map[string]*Stage),
		userBinaryParts: make(map[userBinaryPartID][]byte),
		Semaphore:       make(map[string]*Semaphore),
		SemaphoreIndex:  7,
		discordBot:      config.DiscordBot,
		Name:            config.Name,
		Enable:          config.Enable,
		raviente:        NewRaviente(),
	}

	// Mezeporta
	s.Stages["sl1Ns200p0a0u0"] = NewStage("sl1Ns200p0a0u0")

	// Rasta bar stage
	s.Stages["sl1Ns211p0a0u0"] = NewStage("sl1Ns211p0a0u0")

	// Pallone Carvan
	s.Stages["sl1Ns260p0a0u0"] = NewStage("sl1Ns260p0a0u0")

	// Pallone Guest House 1st Floor
	s.Stages["sl1Ns262p0a0u0"] = NewStage("sl1Ns262p0a0u0")

	// Pallone Guest House 2nd Floor
	s.Stages["sl1Ns263p0a0u0"] = NewStage("sl1Ns263p0a0u0")

	// Diva fountain / prayer fountain.
	s.Stages["sl2Ns379p0a0u0"] = NewStage("sl2Ns379p0a0u0")

	// MezFes
	s.Stages["sl1Ns462p0a0u0"] = NewStage("sl1Ns462p0a0u0")

	return s
}

// Start starts the server in a new goroutine.
func (s *Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}
	s.listener = l

	go s.acceptClients()
	go s.manageSessions()

	// Start the discord bot for chat integration.
	if s.erupeConfig.Discord.Enabled && s.discordBot != nil {
		s.discordBot.Session.AddHandler(s.onDiscordMessage)
	}

	return nil
}

// Shutdown tries to shut down the server gracefully.
func (s *Server) Shutdown() {
	s.Lock()
	s.IsShuttingDown = true
	s.Unlock()

	s.listener.Close()

	close(s.acceptConns)
}

func (s *Server) acceptClients() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.Lock()
			shutdown := s.IsShuttingDown
			s.Unlock()

			if shutdown {
				break
			} else {
				s.logger.Warn("Error accepting client", zap.Error(err))
				continue
			}
		}
		s.acceptConns <- conn
	}
}

func (s *Server) manageSessions() {
	for {
		select {
		case newConn := <-s.acceptConns:
			// Gracefully handle acceptConns channel closing.
			if newConn == nil {
				s.Lock()
				shutdown := s.IsShuttingDown
				s.Unlock()

				if shutdown {
					return
				}
			}

			session := NewSession(s, newConn)

			s.Lock()
			s.Sessions[newConn] = session
			s.Unlock()

			session.Start()

		case delConn := <-s.deleteConns:
			s.Lock()
			delete(s.Sessions, delConn)
			s.Unlock()
		}
	}
}

// BroadcastMHF queues a MHFPacket to be sent to all sessions.
func (s *Server) BroadcastMHF(pkt mhfpacket.MHFPacket, ignoredSession *Session) {
	// Broadcast the data.
	for _, session := range s.Sessions {
		if session == ignoredSession {
			continue
		}

		// Make the header
		bf := byteframe.NewByteFrame()
		bf.WriteUint16(uint16(pkt.Opcode()))

		// Build the packet onto the byteframe.
		// nolint:errcheck // Error return value of `pkt.Build` is not checked
		pkt.Build(bf, session.clientContext)

		// Enqueue in a non-blocking way that drops the packet if the connections send buffer channel is full.
		session.QueueSendNonBlocking(bf.Data())
	}
}

func (s *Server) WorldcastMHF(pkt mhfpacket.MHFPacket, ignoredSession *Session, ignoredChannel *Server) {
	for _, c := range s.Channels {
		if c == ignoredChannel {
			continue
		}
		for _, session := range c.Sessions {
			if session == ignoredSession {
				continue
			}
			bf := byteframe.NewByteFrame()
			bf.WriteUint16(uint16(pkt.Opcode()))
			// nolint:errcheck // Error return value of `pkt.Build` is not checked
			pkt.Build(bf, session.clientContext)
			session.QueueSendNonBlocking(bf.Data())
		}
	}
}

// BroadcastChatMessage broadcasts a simple chat message to all the sessions.
func (s *Server) BroadcastChatMessage(message string) {
	bf := byteframe.NewByteFrame()
	bf.SetLE()
	msgBinChat := &binpacket.MsgBinChat{
		Unk0:       0,
		Type:       5,
		Flags:      0x80,
		Message:    message,
		SenderName: s.Name,
	}
	// nolint:errcheck //  Error return value of `msgBinChat.Build` is not checked
	msgBinChat.Build(bf)

	s.BroadcastMHF(&mhfpacket.MsgSysCastedBinary{
		CharID:         0xFFFFFFFF,
		MessageType:    BinaryMessageTypeChat,
		RawDataPayload: bf.Data(),
	}, nil)
}

func (s *Server) BroadcastRaviente(ip uint32, port uint16, stage []byte, _type uint8) {
	bf := byteframe.NewByteFrame()
	bf.SetLE()
	bf.WriteUint16(0)    // Unk
	bf.WriteUint16(0x43) // Data len
	bf.WriteUint16(3)    // Unk len
	var text string
	switch _type {
	case 2:
		text = "<Great Slaying: Berserk> is being held!"
	case 4:
		text = "<Great Slaying: Extreme> is being held!"
	case 5:
		text = "<Great Slaying: Berserk Practice> is being held!"
	default:
		s.logger.Error("Unk raviente type", zap.Uint8("_type", _type))
	}
	ps.Uint16(bf, text, false)
	bf.WriteBytes([]byte{0x5F, 0x53, 0x00})
	bf.WriteUint32(ip)   // IP address
	bf.WriteUint16(port) // Port
	bf.WriteUint16(0)    // Unk
	bf.WriteBytes(stage)
	s.WorldcastMHF(&mhfpacket.MsgSysCastedBinary{
		CharID:         0x00000000,
		BroadcastType:  BroadcastTypeServer,
		MessageType:    BinaryMessageTypeChat,
		RawDataPayload: bf.Data(),
	}, nil, s)
}

func (s *Server) DiscordChannelSend(charName string, content string) {
	if s.erupeConfig.Discord.Enabled && s.discordBot != nil {
		message := fmt.Sprintf("**%s**: %s", charName, content)
		// nolint:errcheck // rror return value of `s.discordBot.RealtimeChannelSend` is not checked (errcheck)
		s.discordBot.RealtimeChannelSend(message)
	}
}

func (s *Server) FindSessionByCharID(charID uint32) *Session {
	for _, c := range s.Channels {
		for _, client := range c.Sessions {
			if client.CharID == charID {
				return client
			}
		}
	}
	return nil
}

func (s *Server) FindObjectByChar(charID uint32) *Object {
	s.stagesLock.RLock()
	defer s.stagesLock.RUnlock()
	for _, stage := range s.Stages {
		stage.RLock()
		for objId := range stage.Objects {
			obj := stage.Objects[objId]
			if obj.OwnerCharID == charID {
				stage.RUnlock()
				return obj
			}
		}
		stage.RUnlock()
	}

	return nil
}

func (s *Server) NextSemaphoreID() uint32 {
	for {
		exists := false
		s.SemaphoreIndex = s.SemaphoreIndex + 1
		if s.SemaphoreIndex == 0 {
			s.SemaphoreIndex = 7 // Skip reserved indexes
		}
		for _, semaphore := range s.Semaphore {
			if semaphore.Id == s.SemaphoreIndex {
				exists = true
			}
		}
		if !exists {
			break
		}
	}
	return s.SemaphoreIndex
}
