package httpserver

import (
	"erupe-ce/config"
	"erupe-ce/server/channelserver"
	"net/http"
	"strconv"

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

func RunHttpServer(context *HttpServerContext) {
	router := gin.Default()

	router.Use(authMiddleware(context.Token))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	router.GET("/servers", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, context.Servers)
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
