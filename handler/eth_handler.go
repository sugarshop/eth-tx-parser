package handler

import "github.com/gin-gonic/gin"

type ETHHandler struct {

}

// NewETHHandler return ETH handler
func NewETHHandler() *ETHHandler {
	return &ETHHandler{}
}

func (eth *ETHHandler) Register(e *gin.Engine) {
	e.GET("/v1/get_current_block", JSONWrapper(eth.GetCurrentBlock))
	e.POST("/v1/subscribe", JSONWrapper(eth.Subscribe))
	e.GET("/v1/get_transactions", JSONWrapper(eth.GetTransactions))
}

// GetCurrentBlock get last parsed block.
func (eth *ETHHandler) GetCurrentBlock(c *gin.Context) (interface{}, error) {
	return nil, nil
}

// Subscribe subscribe address to server.
func (eth *ETHHandler) Subscribe(c *gin.Context) (interface{}, error) {
	return nil, nil
}

// GetTransactions list of inbound or outbound transactions for an address.
func (eth *ETHHandler) GetTransactions(c *gin.Context) (interface{}, error) {
	return nil, nil
}