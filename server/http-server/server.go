package httpserver

import (
	"encoding/hex"
	"erupe-ce/config"
	"erupe-ce/server/channelserver"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServerContext struct {
	Servers     []*channelserver.Server
	ErupeConfig *config.Config
	Address     string
}

type HttpServerGinContext struct {
	*gin.Context

	Servers     []*channelserver.Server
	ErupeConfig *config.Config
	Address     string
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

func RunHttpServer(context *HttpServerContext) {
	router := gin.Default()

	router.GET("/servers", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, context.Servers)
	})

	router.GET("/raw-binary-data", func(c *gin.Context) {
		data := []RawBinaryKeyEntries{}

		for _, server := range context.Servers {
			for _, stage := range server.Stages {

				binariesKeys := []RawBinaryKeyValue{}

				for key, value := range stage.RawBinaryData {
					fmt.Printf("key: %d,%d", key.Id0, key.Id1)
					fmt.Printf("value [%d bytes]:\n%s\n", len(data), hex.Dump(value))
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

	router.Run(context.Address)
}
