package util

import (
	"context"
	"github.com/gin-gonic/gin"
)

// CtxString type conversion to avoid CI failure
type CtxString string

// RPCContext method to convert gin.Context to context.Context
func RPCContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	return ctx
}