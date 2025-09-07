package ethclient

import (
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/buni/tx-parser/internal/pkg/testing/golden"
	"github.com/stretchr/testify/assert"
)

var (
	update = flag.Bool("update", false, "update .golden files")
	rcpURL = flag.String("rpc-url", "", "core url")
)

func TestGetBlockByNumber(t *testing.T) {
	tests := []struct {
		name         string
		req          GetBlockByNumberRequest
		setupHandler func(t *testing.T, goldenFileName string) http.HandlerFunc
	}{
		{
			name: "get block by number success",
			req: GetBlockByNumberRequest{
				Number: "0x14b66a0",
			},
			setupHandler: golden.HandlerResponse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenFileName := strings.Replace(t.Name(), "/", "_", -1)
			srv := httptest.NewUnstartedServer(nil)
			expect := GetBlockByNumberResponse{}

			rpcURL := *rcpURL
			if rpcURL == "" && !*update {
				srv.Config.Handler = tt.setupHandler(t, goldenFileName)
				srv.Start()
				rpcURL = srv.URL
				t.Cleanup(srv.Close)
			}
			client := NewClient(rpcURL)

			got, err := client.GetBlockByNumber(context.Background(), &tt.req)

			assert.NoError(t, err)
			expect = golden.LoadJSON(t, *got, goldenFileName, update)
			assert.Equal(t, expect, *got)
		})
	}
}

func TestGetCurrentBlock(t *testing.T) {
	tests := []struct {
		name         string
		req          GetCurrentBlockRequest
		setupHandler func(t *testing.T, goldenFileName string) http.HandlerFunc
	}{
		{
			name:         "get block by number success",
			req:          GetCurrentBlockRequest{},
			setupHandler: golden.HandlerResponse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenFileName := strings.Replace(t.Name(), "/", "_", -1)
			srv := httptest.NewUnstartedServer(nil)
			expect := GetCurrentBlockResponse{}

			rpcURL := *rcpURL
			if rpcURL == "" && !*update {
				srv.Config.Handler = tt.setupHandler(t, goldenFileName)
				srv.Start()
				rpcURL = srv.URL
				t.Cleanup(srv.Close)
			}
			client := NewClient(rpcURL)

			got, err := client.GetCurrentBlock(context.Background(), &tt.req)

			assert.NoError(t, err)
			expect = golden.LoadJSON(t, *got, goldenFileName, update)
			assert.Equal(t, expect, *got)
		})
	}
}
