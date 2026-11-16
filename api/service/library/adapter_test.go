package library

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeKnoraDiscovery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "read-key" {
			t.Error("missing source key")
		}
		switch r.URL.Path {
		case "/api/v1/knowledge-bases":
			io.WriteString(w, `{"success":true,"data":[{"id":"kb1","name":"Library"}]}`)
		case "/api/v1/knowledge-bases/kb1":
			io.WriteString(w, `{"success":true,"data":{"id":"kb1","name":"Library"}}`)
		default:
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	client, err := newRemote(server.URL, "read-key")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := client.Collections(context.Background())
	if err != nil || len(rows) != 1 || rows[0].ID != "kb1" {
		t.Fatalf("discovery failed: %v", err)
	}
	row, err := client.Collection(context.Background(), "kb1")
	if err != nil || row.ID != "kb1" {
		t.Fatalf("selection failed: %v", err)
	}
	for _, id := range []string{"../other", "missing"} {
		if _, err := client.Collection(context.Background(), id); !errors.Is(err, ErrNotFound) {
			t.Fatal("invalid selection accepted")
		}
	}
}

func TestWeKnoraRedirectAndURLValidation(t *testing.T) {
	called := false
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	defer target.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, http.StatusFound) }))
	defer server.Close()
	client, err := newRemote(server.URL, "read-key")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Collections(context.Background()); err == nil || called {
		t.Fatal("followed a credential-bearing redirect")
	}
	for _, invalid := range []string{"file:///etc/passwd", "http://user:key@host", "http://169.254.169.254", "http://[::]", "http://localhost?key=secret"} {
		if _, err := ValidateURL(invalid); err == nil {
			t.Errorf("accepted %s", invalid)
		}
	}
}
