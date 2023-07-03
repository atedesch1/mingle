package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/atedesch1/mingle/models"
)

const (
	DefaultRequestQuantity uint = 10
)

type messageIDRequestURI struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

func (h *Handler) getMessage(ctx *gin.Context) {
	var reqURI messageIDRequestURI
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	message, err := h.storage.GetMessage(reqURI.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, message)
}

func (h *Handler) getMessages(ctx *gin.Context) {
	messages, err := h.storage.GetMessages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type getLatestMessagesRequestBody struct {
	Quantity *uint `json:"quantity,omitempty" binding:"min=1"`
}

func (h *Handler) getLatestMessages(ctx *gin.Context) {
	var reqBody getLatestMessagesRequestBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	quantity := DefaultRequestQuantity
	if reqBody.Quantity != nil {
		quantity = *reqBody.Quantity
	}

	messages, err := h.storage.GetLatestMessages(quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type getMessagesRangeRequestBody struct {
	FromID   uint64 `json:"from_id"            binding:"required,min=1"`
	Quantity *uint  `json:"quantity,omitempty" binding:"min=1"`
}

func (h *Handler) getMessagesRange(ctx *gin.Context) {
	var reqBody getMessagesRangeRequestBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	quantity := DefaultRequestQuantity
	if reqBody.Quantity != nil {
		quantity = *reqBody.Quantity
	}

	messages, err := h.storage.GetMessagesRange(reqBody.FromID, quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

func (h *Handler) subscribeToMessages(ctx *gin.Context) {
	// Setup SSE
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	messageChan := make(chan []byte)
	defer close(messageChan)

	go h.storage.SubscribeToMessages(messageChan)

	for {
		select {
		case message := <-messageChan:
			_, _ = ctx.Writer.WriteString("data: " + string(message) + "\n\n")
			flusher.Flush()
		case <-ctx.Request.Context().Done():
			return
		}
	}
}

type createMessageRequestBody struct {
	UserID  uint64 `json:"user_id" binding:"required,min=1"`
	Content string `json:"content" binding:"required"`
}

func (h *Handler) createMessage(ctx *gin.Context) {
	var reqBody createMessageRequestBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := models.MessageCreateParams{
		UserID:  reqBody.UserID,
		Content: reqBody.Content,
	}
	message, err := h.storage.CreateMessage(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, message)
}

func (h *Handler) deleteMessage(ctx *gin.Context) {
	var reqURI messageIDRequestURI
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := h.storage.DeleteMessage(reqURI.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
