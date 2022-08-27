package httpserver

import (
	"erupe-ce/config"
	"erupe-ce/server/channelserver"
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
	router.Run(context.Address)
}
