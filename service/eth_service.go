package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sugarshop/eth-tx-parser/model"
)

// ETHService ETH Transactions data parser service.
type ETHService struct {
	subAddrs []string
	inboundTx map[string]interface{}
	outboundTx map[string]interface{}
}

var _ETHServiceInstance *ETHService

// ETHServiceInstance ETHService singleton
func ETHServiceInstance() *ETHService {
	_ETHServiceInstance = &ETHService{
		subAddrs:   []string{},
		inboundTx:  map[string]interface{}{},
		outboundTx: map[string]interface{}{},
	}
	return _ETHServiceInstance
}

// GetCurrentBlock get current block.
func (s *ETHService) GetCurrentBlock(ctx context.Context) (interface{}, error) {
	num, err := ETHServiceInstance().EthBlockNumber(ctx)
	if err != nil {
		log.Println(ctx, "[GetCurrentBlock]: Error EthBlockNumber, err: ", err)
		return nil, err
	}
	blockInfo, err := ETHServiceInstance().EthGetBlockByNumber(ctx, num)
	if err != nil {
		log.Println(ctx, "[GetCurrentBlock]: Error EthGetBlockByNumber, err: ", err)
		return nil, err
	}
	return blockInfo, nil
}

// Subscribe subscribe an address's inbound/outbound transaction.
func (s *ETHService) Subscribe(ctx context.Context, address string) (interface{}, error) {
	return nil, nil
}

// GetTransactions get address's inbound/outbound transactions
func (s *ETHService) GetTransactions(ctx context.Context, address string) (interface{}, error) {
	return nil, nil
}

// EthBlockNumber returns the number of the most recent block.
func (s *ETHService) EthBlockNumber(ctx context.Context) (string, error) {
	request := &model.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      83, // match response, debug, support multi-request, should be a uniq random number.
	}

	body, err := s.httpJsonRPCPOST(ctx, request)
	if err != nil {
		fmt.Println("[EthBlockNumber]: Error httpJsonRPCPOST request:", err)
		return "", err
	}

	resp := &model.ETHBlockNumberResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		log.Println(ctx, "[EthBlockNumber]: Error Unmarshal, err: ", err)
		return "", err
	}
	hexNumber := resp.Result

	return hexNumber, nil
}

// EthGetBlockByNumber returns information about a block by number.
func (s *ETHService) EthGetBlockByNumber(ctx context.Context, number string) (*model.ETHBlockInfo, error) {
	request := &model.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{number, true},
		ID:      84, // match response, debug, support multi-request, should be a uniq random number.
	}

	body, err := s.httpJsonRPCPOST(ctx, request)
	if err != nil {
		fmt.Println("[EthGetBlockByNumber]: Error httpJsonRPCPOST request:", err)
		return nil, err
	}
	resp := &model.ETHGetBlockByNumberResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		log.Println(ctx, "[EthGetBlockByNumber]: Error Unmarshal, err: ", err)
		return nil, err
	}
	blockInfo := &resp.Result
	return blockInfo, nil
}

func (s *ETHService) httpJsonRPCPOST(ctx context.Context, request *model.JSONRPCRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("[httpJsonRPCPOST]: Error marshaling request:", err)
		return nil, err
	}

	// create HTTP POST request
	url := model.ETHJsonRpcUrl
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("[httpJsonRPCPOST]: Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// HTTP Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[httpJsonRPCPOST]: Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// read resp data.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[httpJsonRPCPOST]: Error reading response:", err)
		return nil, err
	}

	return body, nil
}