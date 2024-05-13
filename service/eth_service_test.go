package service

import (
	"context"
	"github.com/tj/assert"
	"testing"
)

func TestETHService_EthBlockNumber(t *testing.T) {
	ctx := context.Background()
	num, err := ETHServiceInstance().EthBlockNumber(ctx)
	assert.Nil(t, err)
	assert.NotEqual(t, num, "")
}

func TestETHService_EthGetBlockByNumber(t *testing.T) {
	ctx := context.Background()
	num, err := ETHServiceInstance().EthBlockNumber(ctx)
	assert.Nil(t, err)
	assert.NotEqual(t, num, "")
	resp, err := ETHServiceInstance().EthGetBlockByNumber(ctx, num)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestETHService_GetCurrentBlock(t *testing.T) {
	ctx := context.Background()
	blockInfo, err := ETHServiceInstance().GetCurrentBlock(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, blockInfo)
}