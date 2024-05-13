package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sugarshop/eth-tx-parser/service"
	"github.com/sugarshop/eth-tx-parser/util"
	"log"
)

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
	ctx := util.RPCContext(c)
	blockInfo, err := service.ETHServiceInstance().GetCurrentBlock(ctx)
	if err != nil {
		log.Println("[GetCurrentBlock]: GetCurrentBlock err: ", err)
		return nil, err
	}
	return blockInfo, nil
}

// Subscribe subscribe address to server.
func (eth *ETHHandler) Subscribe(c *gin.Context) (interface{}, error) {
	ctx := util.RPCContext(c)
	address := c.Request.Form.Get("address")
	if len(address) == 0 {
		log.Println(ctx, "[Subscribe]: parse address param err")
		return nil, errors.New("parse address param err")
	}
	if err := service.ETHServiceInstance().Subscribe(ctx, address); err != nil {
		log.Println(ctx, "[Subscribe]: Subscribe err: ", err)
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// GetTransactions list of inbound or outbound transactions for an address.
func (eth *ETHHandler) GetTransactions(c *gin.Context) (interface{}, error) {
	ctx := util.RPCContext(c)
	address := c.Request.Form.Get("address")
	if len(address) == 0 {
		log.Println(ctx, "[GetTransactions]: parse address param err")
		return nil, errors.New("parse address param err")
	}
	transactions, err := service.ETHServiceInstance().GetTransactions(ctx, address)
	if err != nil {
		log.Println(ctx, "[GetTransactions]: GetTransactions err: ", err)
		return nil, err
	}
	return map[string]interface{} {
		"transactions": transactions,
	}, nil
}