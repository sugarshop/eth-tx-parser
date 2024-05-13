package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sugarshop/eth-tx-parser/model"
)

// ETHService ETH Transactions data parser service.
type ETHService struct {
	recentBlockNumer int64 // the most recent block number I have ever oberve.
	addrRWMutex sync.RWMutex
	subAddrs map[string]bool
	txRWMutex sync.RWMutex
	transactions map[string][]*model.ETHTransaction
}

var (
	eTHServiceInstance *ETHService
	eTHServiceOnce sync.Once
)

// ETHServiceInstance ETHService singleton
func ETHServiceInstance() *ETHService {
	eTHServiceOnce.Do(func() {
		eTHServiceInstance = &ETHService{
			subAddrs:   map[string]bool{},
			transactions:  map[string][]*model.ETHTransaction{},
		}
		ctx := context.Background()
		dec, err := eTHServiceInstance.ETHBlockDecimalNumber(ctx)
		if err != nil {
			log.Panicln(ctx, "[ETHServiceInstance]: Panic, Error ETHBlockDecimalNumber, err: ", err)
		}
		eTHServiceInstance.recentBlockNumer = dec

		go func() {
			// query eth block number per second.
			// if new block number appear, getBlockByNumber.
			// parse tx into inbount/outbound.
			for range time.Tick(1 * time.Second) {
				if err := eTHServiceInstance.load(ctx); err != nil {
					log.Println(ctx, "[ETHServiceInstance]: eTHServiceInstance load err: %v", err)
				}
			}
		}()
	})

	return eTHServiceInstance
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
func (s *ETHService) Subscribe(ctx context.Context, address string) error {
	s.addrRWMutex.Lock()
	s.subAddrs[address] = true
	s.addrRWMutex.Unlock()
	return nil
}

// GetTransactions get address's inbound/outbound transactions
func (s *ETHService) GetTransactions(ctx context.Context, address string) ([]*model.ETHTransaction, error) {
	s.txRWMutex.RLock()
	transactions, ok := s.transactions[address]
	if !ok {
		transactions = make([]*model.ETHTransaction, 0)
	}
	s.txRWMutex.RUnlock()
	return transactions, nil
}

// load load transactions via address.
func (s *ETHService) load(ctx context.Context) error {
	// 1. query new block number.
	num, err := s.ETHBlockDecimalNumber(ctx)
	if err != nil {
		log.Println(ctx, "[load]: Error EthBlockNumber request:", err)
		return err
	}
	// 2. compare, if no new block, return
	if s.recentBlockNumer >= num {
		// no new block, return.
		return nil
	}
	// 3. update block number.
	s.recentBlockNumer = num
	log.Println(ctx, "[ETHService]: Block Number:", num)
	// 4. parse block transactions.
	err = s.ParseTransactions(ctx, num)
	if err != nil {
		log.Println(ctx, "[load]: Error ParseTransactions request:", err)
		return err
	}
	return nil
}

// ParseTransactions parse block transactions.
func (s *ETHService) ParseTransactions(ctx context.Context, number int64) error {
	hexStr := fmt.Sprintf("0x%x", number)
	blockInfo, err := s.EthGetBlockByNumber(ctx, hexStr)
	if err != nil {
		log.Println(ctx, "[ParseTransactions]: Error EthGetBlockByNumber request:", err)
		return err
	}
	transactions := blockInfo.Transactions
	for _, tx := range transactions {
		// if a key exists in map, store it.
		s.addrRWMutex.RLock()
		s.txRWMutex.Lock()
		if _, ok := s.subAddrs[tx.From]; ok {
			// outboundTx: From -> To
			if txList, okk := s.transactions[tx.From]; okk {
				txList = append(txList, tx)
				s.transactions[tx.From] = txList
			} else {
				s.transactions[tx.From] = []*model.ETHTransaction{tx}
			}
		}
		if _, ok := s.subAddrs[tx.To]; ok {
			// inboundTx: From -> To
			if txList, okk := s.transactions[tx.To]; okk {
				txList = append(txList, tx)
				s.transactions[tx.To] = txList
			} else {
				s.transactions[tx.To] = []*model.ETHTransaction{tx}
			}
		}
		s.addrRWMutex.RUnlock()
		s.txRWMutex.Unlock()
	}
	return nil
}

// ETHBlockDecimalNumber return the decimal number of the most recent block.
func (s *ETHService) ETHBlockDecimalNumber(ctx context.Context) (int64, error) {
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
func (s *ETHService) EthBlockNumber(ctx context.Context) (string, error) {
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
func (s *ETHService) EthGetBlockByNumber(ctx context.Context, number string) (*model.ETHBlockInfo, error) {
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

func (s *ETHService) httpJsonRPCPOST(ctx context.Context, request *model.JSONRPCRequest) ([]byte, error) {
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