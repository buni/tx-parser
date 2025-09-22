package golden

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func Load(t *testing.T, body []byte, goldenFileName string, update *bool) []byte {
	t.Helper()

	goldenFilePath := filepath.Join("testdata", goldenFileName+".golden")

	if *update {
		if body == nil {
			t.Fatalf("no body returned from updateFunc: %s", goldenFileName)
		}

		err := os.WriteFile(goldenFilePath, body, 0o644) //nolint:gosec
		if err != nil {
			t.Fatalf("write golden file: %v", err)
		}
	}

	golden, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatalf("read golden file: %v", err)
		return nil
	}

	return golden
}

func LoadJSON[T any](t *testing.T, body T, goldenFileName string, update *bool) (resp T) {
	t.Helper()

	if *update {
		buf := &bytes.Buffer{}

		enc := json.NewEncoder(buf)
		enc.SetIndent("", "  ")

		err := enc.Encode(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
			return resp
		}
		_ = Load(t, buf.Bytes(), goldenFileName, update)

		return body
	}

	golden := Load(t, nil, goldenFileName, update)

	err := json.Unmarshal(golden, &resp)
	if err != nil {
		t.Fatalf("unmarshal golden file: %v", err)
		return resp
	}
	return resp
}
