package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/atedesch1/mingle/db"
)

type Handler struct {
	addr   string
	router *gin.Engine
	storage  db.Storage
}

func NewHandler(storage db.Storage, addr string) *Handler {
	handler := &Handler{
		addr:  addr,
		storage: storage,
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("v1")
	{
		userGroup := v1.Group("user")
		{
			userGroup.GET("/:id", handler.getUser)
			userGroup.GET("/", handler.getUsers)
			userGroup.POST("/", handler.createUser)
			userGroup.DELETE("/:id", handler.deleteUser)
		}

		messageGroup := v1.Group("message")
		{
			messageGroup.GET("/:id", handler.getMessage)
			messageGroup.GET("/", handler.getMessages)
			messageGroup.GET("/latest", handler.getLatestMessages)
			messageGroup.GET("/range", handler.getMessagesRange)
			messageGroup.POST("/", handler.createMessage)
			messageGroup.DELETE("/:id", handler.deleteMessage)
		}
	}

	handler.router = router
	return handler
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (h *Handler) Serve() error {
	return h.router.Run(h.addr)
}
