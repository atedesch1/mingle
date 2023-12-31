package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"

	. "github.com/atedesch1/mingle/errors"
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
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, message)
}

func (h *Handler) getMessages(ctx *gin.Context) {
	if ctx.Query("expand") == "user" {
		messages, err := h.storage.GetMessagesExpandedUser()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, messages)
	}

	messages, err := h.storage.GetMessages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type getLatestMessagesRequestBody struct {
	Quantity *uint `json:"quantity,omitempty"`
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

	if ctx.Query("expand") == "user" {
		messages, err := h.storage.GetLatestMessagesExpandedUser(quantity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, messages)
	}

	messages, err := h.storage.GetLatestMessages(quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type getMessagesRangeRequestBody struct {
	FromID   uint64 `json:"fromId"             binding:"required,min=1"`
	Quantity *uint  `json:"quantity,omitempty"`
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

	if ctx.Query("expand") == "user" {
		messages, err := h.storage.GetMessagesRangeExpandedUser(reqBody.FromID, quantity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, messages)
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
	unsubscribe := make(chan struct{})
	defer close(unsubscribe)

	go h.storage.SubscribeToMessages(messageChan, unsubscribe)

	marshalfromDBToJSON := func(data []byte, v interface{}) ([]byte, error) {
		err := jsoniter.Config{
			TagKey: "db",
		}.Froze().Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
		json, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return json, nil
	}

	for {
		select {
		case message := <-messageChan:
			var msg models.MessageExpandedUser
			json, err := marshalfromDBToJSON(message, &msg)
			if err != nil {
				unsubscribe <- struct{}{}
				return
			}

			_, _ = ctx.Writer.WriteString("data: " + string(json) + "\n\n")
			flusher.Flush()
		case <-ctx.Request.Context().Done():
			unsubscribe <- struct{}{}
			return
		case <-unsubscribe:
			return
		}
	}
}

type createMessageRequestBody struct {
	UserID  uint64 `json:"userId"  binding:"required,min=1"`
	Content string `json:"content" binding:"required"`
}

func (h *Handler) createMessage(ctx *gin.Context) {
	var reqBody createMessageRequestBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, err := h.storage.GetUser(reqBody.UserID); err != nil {
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
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

	if _, err := h.storage.GetMessage(reqURI.ID); err != nil {
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	if err := h.storage.DeleteMessage(reqURI.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
