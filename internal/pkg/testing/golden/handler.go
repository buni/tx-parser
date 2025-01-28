package golden

import (
	"net/http"
	"testing"

	"github.com/buni/tx-parser/internal/pkg/testing/testutils"
)

func HandlerResponse(t *testing.T, goldenFileName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		respBody := Load(t, nil, goldenFileName, testutils.ToPtr(false))
		w.Write(respBody) //nolint
	})
}

func HandlerResponseWithStatus(status int) func(t *testing.T, goldenFileName string) http.HandlerFunc {
	return func(t *testing.T, goldenFileName string) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			respBody := Load(t, nil, goldenFileName, testutils.ToPtr(false))
			w.WriteHeader(status)
			w.Write(respBody) //nolint
		})
	}
}
