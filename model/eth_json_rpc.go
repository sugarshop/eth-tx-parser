package model

const ETHJsonRpcUrl = "https://eth.public-rpc.com"

// JSONRPCRequest represents the structure of the JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// ETHBlockNumberResponse response of the ethBlockNumber request
type ETHBlockNumberResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

// ETHGetBlockByNumberResponse response of the eth_getBlockByNumber request
type ETHGetBlockByNumberResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  *ETHBlockInfo `json:"result"`
}

type ETHTransaction struct {
	BlockHash            string   `json:"blockHash"`
	BlockNumber          string   `json:"blockNumber"`
	From                 string   `json:"from"`
	Gas                  string   `json:"gas"`
	GasPrice             string   `json:"gasPrice"`
	MaxPriorityFeePerGas string   `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string   `json:"maxFeePerGas"`
	Hash                 string   `json:"hash"`
	Input                string   `json:"input"`
	Nonce                string   `json:"nonce"`
	To                   string   `json:"to"`
	TransactionIndex     string   `json:"transactionIndex"`
	Value                string   `json:"value"`
	Type                 string   `json:"type"`
	AccessList           []interface{} `json:"accessList"`
	ChainID              string   `json:"chainId"`
	V                    string   `json:"v"`
	YParity              string   `json:"yParity"`
	R                    string   `json:"r"`
	S                    string   `json:"s"`
}

type ETHBlockInfo struct {
	BaseFeePerGas         string `json:"baseFeePerGas"`
	BlobGasUsed           string `json:"blobGasUsed"`
	Difficulty            string `json:"difficulty"`
	ExcessBlobGas         string `json:"excessBlobGas"`
	ExtraData             string `json:"extraData"`
	GasLimit              string `json:"gasLimit"`
	GasUsed               string `json:"gasUsed"`
	Hash                  string `json:"hash"`
	LogsBloom             string `json:"logsBloom"`
	Miner                 string `json:"miner"`
	MixHash               string `json:"mixHash"`
	Nonce                 string `json:"nonce"`
	Number                string `json:"number"`
	ParentBeaconBlockRoot string `json:"parentBeaconBlockRoot"`
	ParentHash            string `json:"parentHash"`
	ReceiptsRoot          string `json:"receiptsRoot"`
	Sha3Uncles            string `json:"sha3Uncles"`
	Size                  string `json:"size"`
	StateRoot             string `json:"stateRoot"`
	Timestamp             string `json:"timestamp"`
	TotalDifficulty       string `json:"totalDifficulty"`
	Transactions          []*ETHTransaction `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []interface{} `json:"uncles"`
	Withdrawals      []*ETHWithdraw `json:"withdrawals"`
	WithdrawalsRoot string `json:"withdrawalsRoot"`
}

// ETHWithdraw ETH Withdraw infomation
type ETHWithdraw struct {
	Index          string `json:"index"`
	ValidatorIndex string `json:"validatorIndex"`
	Address        string `json:"address"`
	Amount         string `json:"amount"`
}