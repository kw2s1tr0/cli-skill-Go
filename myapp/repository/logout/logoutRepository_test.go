package logout

import (
	"aiagentcliapp/repository"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeTokenStore struct {
	accessToken string
	getErr      error
	deleteErr   error
	deleted     bool
}

func (store *fakeTokenStore) GetAccessToken(ctx context.Context) (string, error) {
	return store.accessToken, store.getErr
}

func (store *fakeTokenStore) DeleteAccessToken(ctx context.Context) error {
	store.deleted = true
	return store.deleteErr
}

func TestLogout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", request.Method)
		}
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
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	if err := NewRepository(client, tokenStore).Logout(context.Background()); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if !tokenStore.deleted {
		t.Fatal("DeleteAccessToken() was not called")
	}
}

func TestLogoutDoesNotDeleteAccessTokenOnAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	err := NewRepository(client, tokenStore).Logout(context.Background())
	if err == nil {
		t.Fatal("Logout() error = nil, want error")
	}
	if tokenStore.deleted {
		t.Fatal("DeleteAccessToken() was called")
	}
}

func TestLogoutReturnsTokenStoreError(t *testing.T) {
	client := repository.NewClient("http://example.com", http.DefaultClient)
	tokenStore := &fakeTokenStore{getErr: errors.New("missing token")}
	err := NewRepository(client, tokenStore).Logout(context.Background())
	if err == nil {
		t.Fatal("Logout() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "get access token: missing token") {
		t.Fatalf("Logout() error = %q, want token error", err)
	}
}
