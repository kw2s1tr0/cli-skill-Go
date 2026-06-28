package employee

import (
	rootrepository "aiagentcliapp/repository"
	repositoryemployee "aiagentcliapp/repository/employee"
	"aiagentcliapp/repository/employee/output"
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

func TestEmployees(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]output.Output{
			{ID: 1, Name: "Alice", Email: "alice@example.com", Department: "Sales", Position: "Manager"},
			{ID: 2, Name: "Bob", Email: "bob@example.com", Department: "Engineering", Position: "Engineer"},
		})
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	service := NewService(repositoryemployee.NewRepository(client), &fakeTokenStore{accessToken: "1|secret-token"})

	actual, err := service.Employees(context.Background())
	if err != nil {
		t.Fatalf("Employees() error = %v", err)
	}
	if len(actual) != 2 ||
		actual[0].ID != 1 ||
		actual[0].Name != "Alice" ||
		actual[0].Email != "alice@example.com" ||
		actual[0].Department != "Sales" ||
		actual[0].Position != "Manager" ||
		actual[1].ID != 2 ||
		actual[1].Name != "Bob" ||
		actual[1].Email != "bob@example.com" ||
		actual[1].Department != "Engineering" ||
		actual[1].Position != "Engineer" {
		t.Fatalf("Employees() = %#v", actual)
	}
}

func TestEmployeesReturnsTokenStoreError(t *testing.T) {
	service := NewService(
		repositoryemployee.NewRepository(rootrepository.NewClient("http://example.com", http.DefaultClient)),
		&fakeTokenStore{err: errors.New("missing token")},
	)

	_, err := service.Employees(context.Background())
	if err == nil {
		t.Fatal("Employees() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "get access token: missing token") {
		t.Fatalf("Employees() error = %q, want token error", err)
	}
}
