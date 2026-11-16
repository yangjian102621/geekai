package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLibrarySourceNeverSerializesCredential(t *testing.T) {
	data, err := json.Marshal(LibrarySource{Name: "Example", Credential: "secret-api-key"})
	if err != nil || strings.Contains(string(data), "secret") || strings.Contains(string(data), "credential") {
		t.Fatalf("credential leaked: %s (%v)", data, err)
	}
}
