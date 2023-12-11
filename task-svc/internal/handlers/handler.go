package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {

}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET(
		"/ping", func(ctx *gin.Context) {ctx.String(http.StatusOK, "pong")})

	return router
}