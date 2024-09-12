package handler

import (
	"github.com/gin-gonic/gin"
	game "github.com/ntl-y/gameoflife/server"
)

type Handler struct {
	hub *game.Hub
}

func NewHandler(hub *game.Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) ServeWsHandler(c *gin.Context) {
	game.ServeWs(h.hub, c)
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()
	r.GET("/", h.ServeWsHandler)
	return r
}
