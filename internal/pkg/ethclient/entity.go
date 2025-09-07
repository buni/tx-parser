package ethclient

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type HexBigInt struct {
	*big.Int
}

func NewHexBigInt(value *big.Int) *HexBigInt {
	return &HexBigInt{Int: value}
}

func (h *HexBigInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err //nolint:wrapcheck
	}

	s = strings.TrimPrefix(s, "0x")

	h.Int = new(big.Int)
	if _, ok := h.SetString(s, 16); !ok {
		return fmt.Errorf("invalid hex string: %s", s)
	}

	return nil
}

func (h *HexBigInt) MarshalJSON() ([]byte, error) {
	if h.Int == nil {
		return json.Marshal("0x0") //nolint:wrapcheck
	}
	return json.Marshal("0x" + h.Text(16)) //nolint:wrapcheck
}

type ErrorResponse struct {
	StatusCode int
	Body       []byte
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("status code: %d, body: %s", e.StatusCode, string(e.Body))
}

type Request struct {
	ID      int    `json:"id"`
	Method  string `json:"method"`
	JSONRPC string `json:"jsonrpc"`
	Params  []any  `json:"params,omitempty"`
}

type GetBlockByNumberRequest struct {
	Number string `json:"number"`
}

type Response[T any] struct {
	ID      int    `json:"id"`
	Result  *T     `json:"result"`
	JSONRPC string `json:"jsonrpc"`
}

type Transaction struct {
	BlockHash            string     `json:"blockHash"`
	BlockNumber          string     `json:"blockNumber"`
	ChainID              string     `json:"chainId,omitempty"`
	From                 string     `json:"from"`
	Gas                  *HexBigInt `json:"gas"`
	GasPrice             *HexBigInt `json:"gasPrice"`
	Hash                 string     `json:"hash"`
	Input                string     `json:"input"`
	MaxFeePerGas         string     `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string     `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                string     `json:"nonce"`
	R                    string     `json:"r"`
	S                    string     `json:"s"`
	To                   string     `json:"to"`
	TransactionIndex     string     `json:"transactionIndex"`
	Type                 string     `json:"type"`
	V                    string     `json:"v"`
	Value                *HexBigInt `json:"value"`
	YParity              string     `json:"yParity,omitempty"`
	BlobVersionedHashes  []string   `json:"blobVersionedHashes,omitempty"`
	MaxFeePerBlobGas     string     `json:"maxFeePerBlobGas,omitempty"`
}

type Withdrawal struct {
	Address        string `json:"address"`
	Amount         string `json:"amount"`
	Index          string `json:"index"`
	ValidatorIndex string `json:"validatorIndex"`
}

type Block struct {
	BaseFeePerGas         string        `json:"baseFeePerGas"`
	BlobGasUsed           string        `json:"blobGasUsed"`
	Difficulty            string        `json:"difficulty"`
	ExcessBlobGas         string        `json:"excessBlobGas"`
	ExtraData             string        `json:"extraData"`
	GasLimit              string        `json:"gasLimit"`
	GasUsed               string        `json:"gasUsed"`
	Hash                  string        `json:"hash"`
	LogsBloom             string        `json:"logsBloom"`
	Miner                 string        `json:"miner"`
	MixHash               string        `json:"mixHash"`
	Nonce                 string        `json:"nonce"`
	Number                string        `json:"number"`
	ParentBeaconBlockRoot string        `json:"parentBeaconBlockRoot"`
	ParentHash            string        `json:"parentHash"`
	ReceiptsRoot          string        `json:"receiptsRoot"`
	Sha3Uncles            string        `json:"sha3Uncles"`
	Size                  string        `json:"size"`
	StateRoot             string        `json:"stateRoot"`
	Timestamp             string        `json:"timestamp"`
	TotalDifficulty       string        `json:"totalDifficulty"`
	Transactions          []Transaction `json:"transactions"`
	TransactionsRoot      string        `json:"transactionsRoot"`
	Withdrawals           []Withdrawal  `json:"withdrawals"`
	WithdrawalsRoot       string        `json:"withdrawalsRoot"`
}

type GetBlockByNumberResponse struct {
	Response[Block]
}

type GetCurrentBlockRequest struct{}

type GetCurrentBlockResponse struct {
	Response[string]
}
