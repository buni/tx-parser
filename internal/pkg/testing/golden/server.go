package golden

import (
	"net/http/httptest"
	"testing"
)

// NewUnstartedServer returns a new httptest.Server that has not been started.
// The server is automatically closed when the test ends.
func NewUnstartedServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewUnstartedServer(nil)
	t.Cleanup(srv.Close)
	return srv
}
