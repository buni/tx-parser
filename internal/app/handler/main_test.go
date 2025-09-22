package handler_test

import (
	"testing"

	httpin_integration "github.com/ggicci/httpin/integration" //nolint
	"github.com/go-chi/chi/v5"
)

func TestMain(m *testing.M) {
	httpin_integration.UseGochiURLParam("path", chi.URLParam)
	m.Run()
}
