package ethclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	jsonrpcVersion = "2.0"
)

//go:generate go tool mockgen -source=client.go -destination=mock/client_mocks.go -package ethclient_mock

type Client interface {
	GetCurrentBlock(ctx context.Context, _ *GetCurrentBlockRequest) (*GetCurrentBlockResponse, error)
	GetBlockByNumber(ctx context.Context, req *GetBlockByNumberRequest) (*GetBlockByNumberResponse, error)
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClient struct {
	baseURL    string
	httpClient Doer
}

func NewClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL:    baseURL,
		httpClient: cleanhttp.DefaultPooledClient(),
	}
}

func (c *HTTPClient) GetCurrentBlock(ctx context.Context, _ *GetCurrentBlockRequest) (*GetCurrentBlockResponse, error) {
	resp := &GetCurrentBlockResponse{}
	rpcReq := &Request{
		JSONRPC: jsonrpcVersion,
		Method:  "eth_blockNumber",
		ID:      1,
	}

	req, err := createRPCRequest(ctx, c.baseURL, rpcReq)
	if err != nil {
		return nil, fmt.Errorf("create rpc request: %w", err)
	}

	resp.Response, err = doRequest[string](req, c.httpClient)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *HTTPClient) GetBlockByNumber(ctx context.Context, req *GetBlockByNumberRequest) (*GetBlockByNumberResponse, error) {
	resp := &GetBlockByNumberResponse{}
	rpcReq := &Request{
		JSONRPC: jsonrpcVersion,
		Method:  "eth_getBlockByNumber",
		Params:  []any{req.Number, true}, // only hydrated transactions are supported by this method
		ID:      1,
	}

	httpReq, err := createRPCRequest(ctx, c.baseURL, rpcReq)
	if err != nil {
		return nil, fmt.Errorf("create rpc request: %w", err)
	}

	resp.Response, err = doRequest[Block](httpReq, c.httpClient)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func createRPCRequest[T any](ctx context.Context, baseURL string, body T) (*http.Request, error) {
	buff := &bytes.Buffer{}
	enc := json.NewEncoder(buff)

	if err := enc.Encode(body); err != nil {
		return nil, fmt.Errorf("encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, buff)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func doRequest[T any](req *http.Request, client Doer) (result Response[T], err error) {
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	respBuff := &bytes.Buffer{}

	teedReader := io.TeeReader(resp.Body, respBuff)

	var rpcResp Response[T]

	if err = json.NewDecoder(teedReader).Decode(&rpcResp); err != nil {
		return result, fmt.Errorf("decode response: %w", err)
	}

	if resp.StatusCode > 399 {
		return result, &ErrorResponse{
			StatusCode: resp.StatusCode,
			Body:       respBuff.Bytes(),
		}
	}

	if rpcResp.Result == nil {
		return result, &ErrorResponse{
			StatusCode: http.StatusNotFound,
			Body:       respBuff.Bytes(),
		}
	}

	return rpcResp, nil
}
