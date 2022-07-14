package channelserver

import (
	"fmt"
	"net"
	"sync"

	"github.com/Andoryuuta/byteframe"
	"github.com/Solenataris/Erupe/config"
	"github.com/Solenataris/Erupe/network/binpacket"
	"github.com/Solenataris/Erupe/network/mhfpacket"
	"github.com/Solenataris/Erupe/server/discordbot"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type StageIdType = string

const (
	// GlobalStage is the stage that is used for all users.
	MezeportaStageId    StageIdType = "sl1Ns200p0a0u0"
	GuildHallLv1StageId StageIdType = "sl1Ns202p0a0u0"
	GuildHallLv2StageId StageIdType = "sl1Ns203p0a0u0"
	GuildHallLv3StageId StageIdType = "sl1Ns204p0a0u0"
	PugiFarmStageId     StageIdType = "sl1Ns205p0a0u0"
	RastaBarStageId     StageIdType = "sl1Ns211p0a0u0"
	CarvaneStageId      StageIdType = "sl1Ns260p0a0u0"
	GookFarmStageId     StageIdType = "sl1Ns265p0a0u0"
	DivaFountainStageId StageIdType = "sl2Ns379p0a0u0"
	DivaHallStageId     StageIdType = "sl1Ns445p0a0u0"
	MezFesStageId       StageIdType = "sl1Ns462p0a0u0"
)

// Config struct allows configuring the server.
type Config struct {
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
	logger         *zap.Logger
	db             *sqlx.DB
	erupeConfig    *config.Config
	acceptConns    chan net.Conn
	deleteConns    chan net.Conn
	sessions       map[net.Conn]*Session
	listener       net.Listener // Listener that is created when Server.Start is called.
	isShuttingDown bool

	stagesLock sync.RWMutex
	stages     map[string]*Stage

	// UserBinary
	userBinaryPartsLock sync.RWMutex
	userBinaryParts     map[userBinaryPartID][]byte

	// Semaphore
	semaphoreLock sync.RWMutex
	semaphore     map[string]*Semaphore

	// Discord chat integration
	discordBot *discordbot.DiscordBot

	name   string
	enable bool

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
		logger:          config.Logger,
		db:              config.DB,
		erupeConfig:     config.ErupeConfig,
		acceptConns:     make(chan net.Conn),
		deleteConns:     make(chan net.Conn),
		sessions:        make(map[net.Conn]*Session),
		stages:          make(map[string]*Stage),
		userBinaryParts: make(map[userBinaryPartID][]byte),
		semaphore:       make(map[string]*Semaphore),
		discordBot:      config.DiscordBot,
		name:            config.Name,
		enable:          config.Enable,
		raviente:        NewRaviente(),
	}

	// Mezeporta
	s.stages["sl1Ns200p0a0u0"] = NewStage("sl1Ns200p0a0u0")

	// Guild Hall LV1
	s.stages["sl1Ns202p0a0u0"] = NewStage("sl1Ns202p0a0u0")

	// Guild Hall LV2
	s.stages["sl1Ns203p0a0u0"] = NewStage("sl1Ns203p0a0u0")

	// Guild Hall LV3
	s.stages["sl1Ns204p0a0u0"] = NewStage("sl1Ns204p0a0u0")

	// Pugi Farm
	s.stages["sl1Ns205p0a0u0"] = NewStage("sl1Ns205p0a0u0")

	// Rasta bar stage
	s.stages["sl1Ns211p0a0u0"] = NewStage("sl1Ns211p0a0u0")

	// Pallone Carvan
	s.stages["sl1Ns260p0a0u0"] = NewStage("sl1Ns260p0a0u0")

	// Gook Farm
	s.stages["sl1Ns265p0a0u0"] = NewStage("sl1Ns265p0a0u0")

	// Diva fountain / prayer fountain.
	s.stages["sl2Ns379p0a0u0"] = NewStage("sl2Ns379p0a0u0")

	// Diva Hall
	s.stages["sl1Ns445p0a0u0"] = NewStage("sl1Ns445p0a0u0")

	// MezFes
	s.stages["sl1Ns462p0a0u0"] = NewStage("sl1Ns462p0a0u0")

	return s
}

// Start starts the server in a new goroutine.
func (s *Server) Start(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
	s.isShuttingDown = true
	s.Unlock()

	s.listener.Close()

	close(s.acceptConns)
}

func (s *Server) acceptClients() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.Lock()
			shutdown := s.isShuttingDown
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
				shutdown := s.isShuttingDown
				s.Unlock()

				if shutdown {
					return
				}
			}

			session := NewSession(s, newConn)

			s.Lock()
			s.sessions[newConn] = session
			s.Unlock()

			session.Start()

		case delConn := <-s.deleteConns:
			s.Lock()
			delete(s.sessions, delConn)
			s.Unlock()
		}
	}
}

// BroadcastMHF queues a MHFPacket to be sent to all sessions.
func (s *Server) BroadcastMHF(pkt mhfpacket.MHFPacket, ignoredSession *Session) {
	// Broadcast the data.
	for _, session := range s.sessions {
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

// BroadcastChatMessage broadcasts a simple chat message to all the sessions.
func (s *Server) BroadcastChatMessage(message string) {
	bf := byteframe.NewByteFrame()
	bf.SetLE()
	msgBinChat := &binpacket.MsgBinChat{
		Unk0:       0,
		Type:       5,
		Flags:      0x80,
		Message:    message,
		SenderName: s.name,
	}
	msgBinChat.Build(bf)

	s.BroadcastMHF(&mhfpacket.MsgSysCastedBinary{
		CharID:         0xFFFFFFFF,
		MessageType:    BinaryMessageTypeChat,
		RawDataPayload: bf.Data(),
	}, nil)
}

func (s *Server) DiscordChannelSend(charName string, content string) {
	if s.erupeConfig.Discord.Enabled && s.discordBot != nil {
		message := fmt.Sprintf("**%s** : %s", charName, content)
		s.discordBot.RealtimeChannelSend(message)
	}
}

func (s *Server) FindSessionByCharID(charID uint32) *Session {
	s.stagesLock.RLock()
	defer s.stagesLock.RUnlock()
	for _, stage := range s.stages {
		stage.RLock()
		for client := range stage.clients {
			if client.charID == charID {
				stage.RUnlock()
				return client
			}
		}
		stage.RUnlock()
	}

	return nil
}

func (s *Server) FindStageObjectByChar(charID uint32) *StageObject {
	s.stagesLock.RLock()
	defer s.stagesLock.RUnlock()
	for _, stage := range s.stages {
		stage.RLock()
		for objId := range stage.objects {
			obj := stage.objects[objId]
			if obj.ownerCharID == charID {
				stage.RUnlock()
				return obj
			}
		}
		stage.RUnlock()
	}

	return nil
}
