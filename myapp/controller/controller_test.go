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
	deleted     bool
}

func (store *fakeTokenStore) SaveAccessToken(ctx context.Context, accessToken string) error {
	store.accessToken = accessToken
	return nil
}

func (store *fakeTokenStore) GetAccessToken(ctx context.Context) (string, error) {
	return store.accessToken, nil
}

func (store *fakeTokenStore) DeleteAccessToken(ctx context.Context) error {
	store.deleted = true
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

func TestMe(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/auth/me" {
			t.Errorf("path = %s, want /api/auth/me", request.URL.Path)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]any{
			"id":    123,
			"name":  "User Name",
			"email": "user@example.com",
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"me"},
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
	if got := stdout.String(); got != "ID: 123\nName: User Name\nEmail: user@example.com\n" {
		t.Fatalf("stdout = %q, want user output", got)
	}
}

func TestLogout(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/auth/logout" {
			t.Errorf("path = %s, want /api/auth/logout", request.URL.Path)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"logout"},
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
	if !tokenStore.deleted {
		t.Fatal("DeleteAccessToken() was not called")
	}
	if got := stdout.String(); got != "logout succeeded\n" {
		t.Fatalf("stdout = %q, want logout succeeded", got)
	}
}

func TestDepartments(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/departments" {
			t.Errorf("path = %s, want /api/departments", request.URL.Path)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]map[string]any{
			{"id": 1, "name": "Sales"},
			{"id": 2, "name": "Engineering"},
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"departments"},
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
	if got := stdout.String(); got != "ID: 1\nName: Sales\nID: 2\nName: Engineering\n" {
		t.Fatalf("stdout = %q, want departments output", got)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}

func TestDepartmentsReturnsAPIError(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/departments" {
			t.Errorf("path = %s, want /api/departments", request.URL.Path)
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"departments"},
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
	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want empty", got)
	}
	if !strings.Contains(stderr.String(), "departments: departments: API returned 500 Internal Server Error: server error") {
		t.Fatalf("stderr = %q, want departments API error", stderr.String())
	}
}

func TestPositions(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/positions" {
			t.Errorf("path = %s, want /api/positions", request.URL.Path)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]map[string]any{
			{"id": 1, "name": "Manager"},
			{"id": 2, "name": "Engineer"},
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"positions"},
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
	if got := stdout.String(); got != "ID: 1\nName: Manager\nID: 2\nName: Engineer\n" {
		t.Fatalf("stdout = %q, want positions output", got)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}

func TestPositionsReturnsAPIError(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/positions" {
			t.Errorf("path = %s, want /api/positions", request.URL.Path)
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"positions"},
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
	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want empty", got)
	}
	if !strings.Contains(stderr.String(), "positions: positions: API returned 500 Internal Server Error: server error") {
		t.Fatalf("stderr = %q, want positions API error", stderr.String())
	}
}
