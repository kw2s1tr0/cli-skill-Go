package logout

import (
	rootrepository "aiagentcliapp/repository"
	repositorylogout "aiagentcliapp/repository/logout"
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
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	service := NewService(repositorylogout.NewRepository(client, tokenStore))

	if err := service.Logout(context.Background()); err != nil {
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

	client := rootrepository.NewClient(server.URL, server.Client())
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	service := NewService(repositorylogout.NewRepository(client, tokenStore))

	err := service.Logout(context.Background())
	if err == nil {
		t.Fatal("Logout() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "logout: logout: API returned 401 Unauthorized") {
		t.Fatalf("Logout() error = %q, want wrapped API error", err)
	}
	if tokenStore.deleted {
		t.Fatal("DeleteAccessToken() was called")
	}
}

func TestLogoutReturnsTokenStoreError(t *testing.T) {
	service := NewService(
		repositorylogout.NewRepository(
			rootrepository.NewClient("http://example.com", http.DefaultClient),
			&fakeTokenStore{getErr: errors.New("missing token")},
		),
	)

	err := service.Logout(context.Background())
	if err == nil {
		t.Fatal("Logout() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "get access token: missing token") {
		t.Fatalf("Logout() error = %q, want token error", err)
	}
}
