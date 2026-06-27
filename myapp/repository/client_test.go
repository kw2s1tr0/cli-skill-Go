package repository

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDoJSONFormatsJSONErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(`{"message":"メールアドレスまたはパスワードが正しくありません。"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	err := client.DoJSON(context.Background(), http.MethodGet, "/login", nil, nil)

	if err == nil {
		t.Fatal("DoJSON() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "API returned 401 Unauthorized: メールアドレスまたはパスワードが正しくありません。") {
		t.Fatalf("DoJSON() error = %q, want decoded message", err)
	}
}

func TestDoJSONFormatsEscapedJSONErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(`{"message":"\u30e1\u30fc\u30eb\u30a2\u30c9\u30ec\u30b9\u307e\u305f\u306f\u30d1\u30b9\u30ef\u30fc\u30c9\u304c\u6b63\u3057\u304f\u3042\u308a\u307e\u305b\u3093\u3002"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	err := client.DoJSON(context.Background(), http.MethodGet, "/login", nil, nil)

	if err == nil {
		t.Fatal("DoJSON() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "API returned 401 Unauthorized: メールアドレスまたはパスワードが正しくありません。") {
		t.Fatalf("DoJSON() error = %q, want decoded message", err)
	}
}

func TestDoJSONKeepsPlainTextErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		http.Error(writer, "invalid credentials", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	err := client.DoJSON(context.Background(), http.MethodGet, "/login", nil, nil)

	if err == nil {
		t.Fatal("DoJSON() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "API returned 401 Unauthorized: invalid credentials") {
		t.Fatalf("DoJSON() error = %q, want plain text message", err)
	}
}

func TestDoJSONFormatsEmptyErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	err := client.DoJSON(context.Background(), http.MethodGet, "/login", nil, nil)

	if err == nil {
		t.Fatal("DoJSON() error = nil, want error")
	}
	if err.Error() != "API returned 401 Unauthorized" {
		t.Fatalf("DoJSON() error = %q, want status only", err)
	}
}
