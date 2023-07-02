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


	handler.router = router
	return handler
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (h *Handler) Serve() error {
	return h.router.Run(h.addr)
}
