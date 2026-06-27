package login

import (
	rootrepository "aiagentcliapp/repository"
	repositorylogin "aiagentcliapp/repository/login"
	loginrequest "aiagentcliapp/repository/login/request"
	"aiagentcliapp/service/login/input"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeTokenStore struct {
	accessToken string
	called      bool
	err         error
}

func (store *fakeTokenStore) SaveAccessToken(ctx context.Context, accessToken string) error {
	store.called = true
	store.accessToken = accessToken
	return store.err
}

func TestLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var actual loginrequest.Request
		if err := json.NewDecoder(request.Body).Decode(&actual); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if actual.Email != "user@example.com" || actual.Password != "password" || actual.TokenName != "agent-cli" {
			t.Fatalf("request = %#v", actual)
		}

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"access_token": "1|secret-token",
			"token_type": "Bearer",
			"expires_at": "2026-07-20T12:00:00+00:00"
		}`))
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	tokenStore := &fakeTokenStore{}
	service := NewService(repositorylogin.NewRepository(client), tokenStore)
	loginInput := input.Input{
		Email:     "user@example.com",
		Password:  "password",
		TokenName: "agent-cli",
	}

	if err := service.Login(context.Background(), loginInput); err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if tokenStore.accessToken != "1|secret-token" {
		t.Fatalf("saved access token = %q, want 1|secret-token", tokenStore.accessToken)
	}
}

func TestLoginDoesNotSaveAccessTokenOnAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "invalid credentials", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	tokenStore := &fakeTokenStore{}
	service := NewService(repositorylogin.NewRepository(client), tokenStore)

	err := service.Login(context.Background(), input.Input{})
	if err == nil {
		t.Fatal("Login() error = nil, want error")
	}
	if tokenStore.called {
		t.Fatal("SaveAccessToken() was called")
	}
}

func TestLoginReturnsSaveAccessTokenError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"access_token": "1|secret-token",
			"token_type": "Bearer",
			"expires_at": "2026-07-20T12:00:00+00:00"
		}`))
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	tokenStore := &fakeTokenStore{err: errors.New("keyring unavailable")}
	service := NewService(repositorylogin.NewRepository(client), tokenStore)

	err := service.Login(context.Background(), input.Input{})
	if err == nil {
		t.Fatal("Login() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "save access token: keyring unavailable") {
		t.Fatalf("Login() error = %q, want save access token error", err)
	}
}
