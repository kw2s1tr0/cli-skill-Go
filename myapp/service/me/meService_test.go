package me

import (
	rootrepository "aiagentcliapp/repository"
	repositoryme "aiagentcliapp/repository/me"
	"aiagentcliapp/repository/me/response"
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
	err         error
}

func (store *fakeTokenStore) GetAccessToken(ctx context.Context) (string, error) {
	return store.accessToken, store.err
}

func TestMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response.Response{
			ID:    123,
			Name:  "User Name",
			Email: "user@example.com",
		})
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	service := NewService(repositoryme.NewRepository(client), &fakeTokenStore{accessToken: "1|secret-token"})

	actual, err := service.Me(context.Background())
	if err != nil {
		t.Fatalf("Me() error = %v", err)
	}
	if actual.ID != 123 || actual.Name != "User Name" || actual.Email != "user@example.com" {
		t.Fatalf("Me() = %#v", actual)
	}
}

func TestMeReturnsTokenStoreError(t *testing.T) {
	service := NewService(
		repositoryme.NewRepository(rootrepository.NewClient("http://example.com", http.DefaultClient)),
		&fakeTokenStore{err: errors.New("missing token")},
	)

	_, err := service.Me(context.Background())
	if err == nil {
		t.Fatal("Me() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "get access token: missing token") {
		t.Fatalf("Me() error = %q, want token error", err)
	}
}
