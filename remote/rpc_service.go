package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/sugarshop/eth-tx-parser/model"
)

// RPCService ETH RPC service.
type RPCService struct {}

var (
	rPCServiceInstance *RPCService
	rPCServiceOnce sync.Once
)

// RPCServiceInstance RPCService singleton
func RPCServiceInstance() *RPCService {
	rPCServiceOnce.Do(func() {
		rPCServiceInstance = &RPCService{}
	})

	return rPCServiceInstance
}

// ETHBlockDecimalNumber return the decimal number of the most recent block.
func (s *RPCService) ETHBlockDecimalNumber(ctx context.Context) (int64, error) {
	hexStr, err := s.EthBlockNumber(ctx)
	if err != nil {
		log.Println(ctx, "[ETHBlockDecimalNumber]: Error EthBlockNumber request:", err)
		return 0, err
	}
	if len(hexStr) == 0 {
		log.Println(ctx, "[ETHBlockDecimalNumber]: Error EthBlockNumber request, hexStr length is 0")
		return 0, errors.New("hexStr length is 0")
	}
	// Convert hexadecimal string to decimal integer
	dec, err := strconv.ParseInt(hexStr[2:], 16, 64)
	if err != nil {
		log.Println(ctx, "[ETHBlockDecimalNumber]: Error ParseInt, err: ", err)
		return 0, err
	}
	return dec, nil
}

// EthBlockNumber returns the number of the most recent block.
func (s *RPCService) EthBlockNumber(ctx context.Context) (string, error) {
	request := &model.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      83, // match response, debug, support multi-request, should be a uniq random number.
	}

	body, err := s.httpJsonRPCPOST(ctx, request)
	if err != nil {
		log.Println(ctx, "[EthBlockNumber]: Error httpJsonRPCPOST request:", err)
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
func (s *RPCService) EthGetBlockByNumber(ctx context.Context, number string) (*model.ETHBlockInfo, error) {
	request := &model.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{number, true},
		ID:      84, // match response, debug, support multi-request, should be a uniq random number.
	}

	body, err := s.httpJsonRPCPOST(ctx, request)
	if err != nil {
		log.Println(ctx, "[EthGetBlockByNumber]: Error httpJsonRPCPOST request:", err)
		return nil, err
	}
	resp := &model.ETHGetBlockByNumberResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		log.Println(ctx, "[EthGetBlockByNumber]: Error Unmarshal, err: ", err)
		return nil, err
	}
	blockInfo := resp.Result
	// TODO: if jsonrpc return nil result, retry it.
	if blockInfo == nil {
		log.Println(ctx, "[EthGetBlockByNumber]: empty blockInfo, should retry")
		return nil, errors.New("empty blockInfo")
	}
	return blockInfo, nil
}

func (s *RPCService) httpJsonRPCPOST(ctx context.Context, request *model.JSONRPCRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Println(ctx, "[httpJsonRPCPOST]: Error marshaling request:", err)
		return nil, err
	}

	// create HTTP POST request
	url := model.ETHJsonRpcUrl
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(ctx, "[httpJsonRPCPOST]: Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// HTTP Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(ctx, "[httpJsonRPCPOST]: Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// read resp data.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(ctx, "[httpJsonRPCPOST]: Error reading response:", err)
		return nil, err
	}

	return body, nil
}
