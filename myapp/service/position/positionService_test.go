package position

import (
	rootrepository "aiagentcliapp/repository"
	repositoryposition "aiagentcliapp/repository/position"
	"aiagentcliapp/repository/position/output"
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

func TestPositions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]output.Output{
			{ID: 1, Name: "Manager"},
			{ID: 2, Name: "Engineer"},
		})
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	service := NewService(repositoryposition.NewRepository(client), &fakeTokenStore{accessToken: "1|secret-token"})

	actual, err := service.Positions(context.Background())
	if err != nil {
		t.Fatalf("Positions() error = %v", err)
	}
	if len(actual) != 2 || actual[0].ID != 1 || actual[0].Name != "Manager" || actual[1].ID != 2 || actual[1].Name != "Engineer" {
		t.Fatalf("Positions() = %#v", actual)
	}
}

func TestPositionsReturnsTokenStoreError(t *testing.T) {
	service := NewService(
		repositoryposition.NewRepository(rootrepository.NewClient("http://example.com", http.DefaultClient)),
		&fakeTokenStore{err: errors.New("missing token")},
	)

	_, err := service.Positions(context.Background())
	if err == nil {
		t.Fatal("Positions() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "get access token: missing token") {
		t.Fatalf("Positions() error = %q, want token error", err)
	}
}
