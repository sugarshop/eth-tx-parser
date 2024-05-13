package service

import (
	"context"
	"fmt"
	"github.com/tj/assert"
	"strings"
	"testing"
)

func TestETHService_GetCurrentBlock(t *testing.T) {
	ctx := context.Background()
	instance := ETHServiceInstance()
	blockInfo, err := instance.GetCurrentBlock(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, blockInfo)
	assert.Equal(t, fmt.Sprintf("0x%x", instance.recentBlockNumer), blockInfo.Number)
}

func TestETHService_Subscribe(t *testing.T) {
	ctx := context.Background()
	address := "0x76759058b7a242A86a0367729FAe98803d86891B"
	ETHServiceInstance().Subscribe(ctx, address)
	_, ok := ETHServiceInstance().subAddrs[strings.ToLower(address)]
	assert.Equal(t, ok, true)
}

func TestETHService_GetTransactions(t *testing.T) {
	ctx := context.Background()
	blockNumber := int64(19862630)
	instance := ETHServiceInstance()
	txCaseList := []struct{
		Addr string
		txNum int
	}{
		{"0xae2fc483527b8ef99eb5d9b44875f005ba1fae13", 2},
		{"0x6b75d8af000000e20b7a7ddf000ba900b4009a80", 2},
		{"0x107fe4e8248ae91651668666e82752890d700eec", 1},
		{"0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad", 8},
		{"0x23ae0460537009106915e962fa95dced15479427", 1},
		{"0x6131b5fae19ea4f9d964eac0408e4408b66337b5", 1},
	}
	for _, txCase := range txCaseList {
		instance.Subscribe(ctx, txCase.Addr)
	}
	err := instance.ParseTransactions(ctx, blockNumber)
	assert.Nil(t, err)

	for _, txCase := range txCaseList {
		list, err := instance.GetTransactions(ctx, txCase.Addr)
		assert.Nil(t, err)
		assert.Equal(t, len(list), txCase.txNum)
		for _, tx := range list {
			assert.Condition(t, func() (success bool) {
				return (strings.EqualFold(tx.From, txCase.Addr)) || (strings.EqualFold(tx.To, txCase.Addr))
			})
		}
	}
}