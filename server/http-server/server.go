package httpserver

import (
	"encoding/hex"
	"erupe-ce/config"
	"erupe-ce/server/channelserver"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func RemoveIndex(s []album, index int) []album {
	return append(s[:index], s[index+1:]...)
}
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON((&newAlbum)); err != nil {
		return
	}

	albums = append(albums, newAlbum)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, album := range albums {
		if album.ID == id {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for idx, album := range albums {
		if album.ID == id {
			albums = RemoveIndex(albums, idx)
			c.JSON(http.StatusOK, nil)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

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

	albumRoutes := router.Group("/albums")
	{
		albumRoutes.GET("/", getAlbums)
		albumRoutes.GET("/:id", getAlbumByID)
		albumRoutes.POST("/", postAlbums)
		albumRoutes.DELETE("/:id", deleteAlbumByID)
	}

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
