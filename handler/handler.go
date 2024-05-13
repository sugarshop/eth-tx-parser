package handler

import (
	"github.com/gin-gonic/gin"
)

// Handler Handler
type Handler interface {
	Register(*gin.Engine)
}

func handlers() []Handler {
	return []Handler{
		NewETHHandler(),
	}
}

// Register Register all API endpoints
func Register(e *gin.Engine) {
	for _, h := range handlers() {
		h.Register(e)
	}
}