package remote

import (
	"context"
	"github.com/tj/assert"
	"strings"
	"testing"
)

func TestRPCService_EthBlockNumber(t *testing.T) {
	ctx := context.Background()
	hexStr, err := RPCServiceInstance().EthBlockNumber(ctx)
	assert.Nil(t, err)
	assert.NotEqual(t, hexStr, "")
	assert.Condition(t, func() (success bool) {
		return strings.HasPrefix(hexStr, "0x")
	})
}

func TestRPCService_EthGetBlockByNumber(t *testing.T) {
	ctx := context.Background()
	hexStr, err := RPCServiceInstance().EthBlockNumber(ctx)
	assert.Nil(t, err)
	assert.NotEqual(t, hexStr, "")
	resp, err := RPCServiceInstance().EthGetBlockByNumber(ctx, hexStr)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Number, hexStr)
}

func TestRPCService_ETHBlockDecimalNumber(t *testing.T) {
	ctx := context.Background()
	hexStr, err := RPCServiceInstance().ETHBlockDecimalNumber(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, hexStr, 0)
}
