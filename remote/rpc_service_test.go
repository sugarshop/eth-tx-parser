package remote

import (
	"context"
	"github.com/tj/assert"
	"testing"
)

func TestRPCService_EthBlockNumber(t *testing.T) {
	ctx := context.Background()
	num, err := RPCServiceInstance().EthBlockNumber(ctx)
	assert.Nil(t, err)
	assert.NotEqual(t, num, "")
}

func TestRPCService_EthGetBlockByNumber(t *testing.T) {
	ctx := context.Background()
	num, err := RPCServiceInstance().EthBlockNumber(ctx)
	assert.Nil(t, err)
	assert.NotEqual(t, num, "")
	resp, err := RPCServiceInstance().EthGetBlockByNumber(ctx, num)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestRPCService_ETHBlockDecimalNumber(t *testing.T) {
	ctx := context.Background()
	num, err := RPCServiceInstance().ETHBlockDecimalNumber(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, num, 0)
}
