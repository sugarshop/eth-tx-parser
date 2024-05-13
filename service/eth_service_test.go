package service

import (
	"context"
	"github.com/tj/assert"
	"testing"
)

func TestETHService_GetCurrentBlock(t *testing.T) {
	ctx := context.Background()
	blockInfo, err := ETHServiceInstance().GetCurrentBlock(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, blockInfo)
}