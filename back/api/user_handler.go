package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/atedesch1/mingle/errors"
	"github.com/atedesch1/mingle/models"
)

type userIDRequestURI struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

func (h *Handler) getUser(ctx *gin.Context) {
	var reqURI userIDRequestURI
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := h.storage.GetUser(reqURI.ID)
	if err != nil {
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) getUsers(ctx *gin.Context) {
	users, err := h.storage.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type createUserRequestBody struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) createUser(ctx *gin.Context) {
	var reqBody createUserRequestBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := models.UserCreateParams{
		Name: reqBody.Name,
	}
	user, err := h.storage.CreateUser(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) deleteUser(ctx *gin.Context) {
	var reqURI userIDRequestURI
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, err := h.storage.GetUser(reqURI.ID); err != nil {
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	if err := h.storage.DeleteUser(reqURI.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
