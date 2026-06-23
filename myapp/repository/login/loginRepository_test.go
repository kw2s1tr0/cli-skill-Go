package login

import (
	rootrepository "aiagentcliapp/repository"
	loginrequest "aiagentcliapp/repository/login/request"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	expiresAt := "2026-07-20T12:00:00+00:00"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", request.Method)
		}
		if request.URL.Path != "/api/auth/login" {
			t.Errorf("path = %s, want /api/auth/login", request.URL.Path)
		}
		if accept := request.Header.Get("Accept"); accept != "application/json" {
			t.Errorf("Accept = %q, want application/json", accept)
		}
		if contentType := request.Header.Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", contentType)
		}

		var input loginrequest.Request
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if input.Email != "user@example.com" || input.Password != "password" || input.TokenName != "agent-cli" {
			t.Fatalf("input = %#v", input)
		}

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"access_token": "1|secret-token",
			"token_type": "Bearer",
			"expires_at": "` + expiresAt + `"
		}`))
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	actual, err := NewRepository(client).Login(context.Background(), loginrequest.Request{
		Email:     "user@example.com",
		Password:  "password",
		TokenName: "agent-cli",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if actual.AccessToken != "1|secret-token" || actual.TokenType != "Bearer" {
		t.Fatalf("Login() = %#v", actual)
	}

	wantExpiresAt, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		t.Fatalf("parse expiresAt: %v", err)
	}
	if !actual.ExpiresAt.Equal(wantExpiresAt) {
		t.Fatalf("Login().ExpiresAt = %v, want %v", actual.ExpiresAt, wantExpiresAt)
	}
}

func TestLoginAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		http.Error(writer, "invalid credentials", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	_, err := NewRepository(client).Login(context.Background(), loginrequest.Request{
		Email:     "user@example.com",
		Password:  "wrong",
		TokenName: "agent-cli",
	})
	if err == nil || !strings.Contains(err.Error(), "401 Unauthorized") {
		t.Fatalf("Login() error = %v, want 401 error", err)
	}
}
