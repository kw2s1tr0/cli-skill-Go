package controller

import (
	"aiagentcliapp/repository"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type fakeTokenStore struct {
	accessToken string
}

func (store *fakeTokenStore) SaveAccessToken(ctx context.Context, accessToken string) error {
	store.accessToken = accessToken
	return nil
}

func TestLoginPromptsForPassword(t *testing.T) {
	tokenStore := &fakeTokenStore{}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	var actual struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		TokenName string `json:"token_name"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/auth/login" {
			t.Errorf("path = %s, want /api/auth/login", request.URL.Path)
		}
		if err := json.NewDecoder(request.Body).Decode(&actual); err != nil {
			t.Errorf("Decode() error = %v", err)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]any{
			"access_token": "token",
			"token_type":   "Bearer",
			"expires_at":   time.Now().UTC(),
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	passwordReaderCalled := false
	passwordReader := func(fd int) ([]byte, error) {
		passwordReaderCalled = true
		return []byte("secret"), nil
	}
	var stdout, stderr bytes.Buffer

	code := Run(
		[]string{"login", "--email", "user@example.com", "--token-name", "agent-cli"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitOK {
		t.Fatalf("Run() = %d, want %d; stderr = %q", code, ExitOK, stderr.String())
	}
	if !passwordReaderCalled {
		t.Fatal("ReadPassword() was not called")
	}
	if actual.Email != "user@example.com" || actual.Password != "secret" || actual.TokenName != "agent-cli" {
		t.Fatalf("request = %#v", actual)
	}
	if tokenStore.accessToken != "token" {
		t.Fatalf("saved access token = %q, want token", tokenStore.accessToken)
	}
	if got := stdout.String(); got != "login succeeded\n" {
		t.Fatalf("stdout = %q, want login succeeded", got)
	}
	if got := stderr.String(); got != "Password: \n" {
		t.Fatalf("stderr = %q, want password prompt", got)
	}
}

func TestLoginPasswordReaderFailure(t *testing.T) {
	client := repository.NewClient("http://example.com", http.DefaultClient)
	passwordReader := func(fd int) ([]byte, error) {
		return nil, errors.New("terminal unavailable")
	}
	var stdout, stderr bytes.Buffer

	code := Run(
		[]string{"login", "--email", "user@example.com"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitRuntime {
		t.Fatalf("Run() = %d, want %d", code, ExitRuntime)
	}
	if !strings.Contains(stderr.String(), "read password: terminal unavailable") {
		t.Fatalf("stderr = %q, want reader error", stderr.String())
	}
}

func TestLoginRejectsPasswordFlagAfterReadingPassword(t *testing.T) {
	client := repository.NewClient("http://example.com", http.DefaultClient)
	passwordReaderCalled := false
	passwordReader := func(fd int) ([]byte, error) {
		passwordReaderCalled = true
		return []byte("secret"), nil
	}
	var stdout, stderr bytes.Buffer

	code := Run(
		[]string{"login", "--email", "user@example.com", "--password", "secret"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitUsage {
		t.Fatalf("Run() = %d, want %d", code, ExitUsage)
	}
	if !passwordReaderCalled {
		t.Fatal("ReadPassword() was not called")
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("stderr = %q, want password flag error", stderr.String())
	}
}
