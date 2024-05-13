package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sugarshop/eth-tx-parser/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	engine := gin.New()
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// register other api
	handler.Register(engine)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	// listen and serve on 0.0.0.0:8080
	go func() {
		engine.Run()
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGIN'T
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}