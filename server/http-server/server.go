package httpserver

import (
	"erupe-ce/config"
	"erupe-ce/server/channelserver"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpServerContext struct {
	Servers     []*channelserver.Server
	ErupeConfig *config.Config
	Address     string
	Token       string
}

type HttpServerGinContext struct {
	*gin.Context

	Servers     []*channelserver.Server
	ErupeConfig *config.Config
	Address     string
	Token       string
}

type RawBinaryKeyValue struct {
	Key   channelserver.StageBinaryKey `json:"key"`
	Value []byte                       `json:"value"`
}
type RawBinaryKeyEntries struct {
	ServerId        uint16              `json:"serverId"`
	StageId         string              `json:"stageId"`
	StageBinaryKeys []RawBinaryKeyValue `json:"stageBinaryKeys"`
}

type SessionJSON struct {
	Session  *channelserver.Session `json:"session"`
	ServerId uint16                 `json:"serverId"`
}

type QuestJSON struct {
	StageId             string    `json:"stageId"`
	ReservedClientSlots uint      `json:"reservedClientSlots"`
	MaxPlayers          uint16    `json:"maxPlayers"`
	QuestFilename       string    `json:"questFilename"`
	HasDeparted         bool      `json:"hasDeparted"`
	CreatedAt           time.Time `json:"createdAt"`
	Chars               []uint32  `json:"chars"`
}

func authMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("access_token")
		if auth != token {
			c.JSON(401, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

type NewMessage struct {
	Content    string `json:"content" binding:"required"`
	SenderName string `json:"sender_name" binding:"required"`
}

func RunHttpServer(context *HttpServerContext) {
	router := gin.Default()

	router.Use(authMiddleware(context.Token))

	router.GET("/current-chars", func(c *gin.Context) {
		charIds := []uint32{}

		for _, server := range context.Servers {
			for _, session := range server.Sessions {
				charIds = append(charIds, session.CharID)
			}
		}

		c.JSON(http.StatusOK, charIds)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	router.GET("/servers", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, context.Servers)
	})

	router.POST("/send-message", func(c *gin.Context) {
		var message NewMessage

		if err := c.BindJSON((&message)); err != nil {
			return
		}

		for _, server := range context.Servers {
			server.BroadcastChatMessageWithCustomName(message.SenderName, message.Content)
		}

		c.JSON(http.StatusOK, nil)
	})

	router.POST("/send-message/:char_id", func(c *gin.Context) {
		id64, err := strconv.ParseUint(c.Param("char_id"), 10, 32)

		if err != nil {
			c.JSON(400, gin.H{"error": "ID SHOULD BE NUMBER"})
			return
		}

		id := uint32(id64)

		var message NewMessage

		if err := c.BindJSON((&message)); err != nil {
			return
		}

		for _, server := range context.Servers {
			for _, session := range server.Sessions {
				if session.CharID == id {
					session.SendMessage(message.SenderName, message.Content)
					c.JSON(http.StatusOK, nil)
					return
				}

			}
		}

		c.JSON(404, gin.H{"error": "NOT FOUND"})
	})

	router.GET("/raw-binary-data", func(c *gin.Context) {
		data := []RawBinaryKeyEntries{}

		for _, server := range context.Servers {
			for _, stage := range server.Stages {

				binariesKeys := []RawBinaryKeyValue{}

				for key, value := range stage.RawBinaryData {
					binariesKeys = append(binariesKeys, RawBinaryKeyValue{
						Key:   key,
						Value: value,
					})
				}

				binaryData := RawBinaryKeyEntries{
					ServerId:        server.ID,
					StageId:         stage.Id,
					StageBinaryKeys: binariesKeys,
				}

				data = append(data, binaryData)

			}
		}

		c.IndentedJSON(http.StatusOK, data)
	})

	router.GET("/sessions", func(c *gin.Context) {
		data := []SessionJSON{}

		for _, server := range context.Servers {
			for _, session := range server.Sessions {
				data = append(data, SessionJSON{
					Session:  session,
					ServerId: server.ID,
				})
			}
		}

		c.IndentedJSON(http.StatusOK, data)
	})

	router.GET("/quests", func(c *gin.Context) {
		quests := []QuestJSON{}

		for _, server := range context.Servers {
			for _, stage := range server.Stages {
				if !stage.IsQuest() {
					continue
				}

				chars := []uint32{}

				for charid := range stage.ReservedClientSlots {
					chars = append(chars, charid)
				}

				quest := QuestJSON{
					StageId:             stage.Id,
					ReservedClientSlots: uint(len(stage.ReservedClientSlots)),
					MaxPlayers:          stage.MaxPlayers,
					QuestFilename:       stage.QuestFilename,
					HasDeparted:         len(stage.Sessions) > 0,
					CreatedAt:           stage.CreatedAt,
					Chars:               chars,
				}

				quests = append(quests, quest)
			}
		}

		c.IndentedJSON(http.StatusOK, quests)
	})

	router.POST("/chars/:char_id/disconnect", func(c *gin.Context) {
		id64, err := strconv.ParseUint(c.Param("char_id"), 10, 32)

		if err != nil {
			c.JSON(400, gin.H{"error": "ID SHOULD BE NUMBER"})
			return
		}

		id := uint32(id64)

		for _, server := range context.Servers {
			for _, session := range server.Sessions {
				if session.CharID == id {
					channelserver.LogoutPlayer(session)
					c.JSON(http.StatusOK, nil)
					return
				}

			}
		}

		c.JSON(404, gin.H{"error": "NOT FOUND"})
	})

	router.DELETE("/stages/:stage_id", func(c *gin.Context) {
		stageId := c.Param("stage_id")

		for _, server := range context.Servers {
			if server.Stages[stageId] != nil {
				delete(server.Stages, stageId)
				c.JSON(http.StatusOK, nil)
				return
			}
		}

		c.JSON(404, gin.H{"error": "NOT FOUND"})
	})

	// nolint:errcheck
	router.Run(context.Address)
}
