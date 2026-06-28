package department

import (
	rootrepository "aiagentcliapp/repository"
	repositorydepartment "aiagentcliapp/repository/department"
	"aiagentcliapp/repository/department/output"
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

func TestDepartments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]output.Output{
			{ID: 1, Name: "Sales"},
			{ID: 2, Name: "Engineering"},
		})
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	service := NewService(repositorydepartment.NewRepository(client), &fakeTokenStore{accessToken: "1|secret-token"})

	actual, err := service.Departments(context.Background())
	if err != nil {
		t.Fatalf("Departments() error = %v", err)
	}
	if len(actual) != 2 || actual[0].ID != 1 || actual[0].Name != "Sales" || actual[1].ID != 2 || actual[1].Name != "Engineering" {
		t.Fatalf("Departments() = %#v", actual)
	}
}

func TestDepartmentsReturnsTokenStoreError(t *testing.T) {
	service := NewService(
		repositorydepartment.NewRepository(rootrepository.NewClient("http://example.com", http.DefaultClient)),
		&fakeTokenStore{err: errors.New("missing token")},
	)

	_, err := service.Departments(context.Background())
	if err == nil {
		t.Fatal("Departments() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "get access token: missing token") {
		t.Fatalf("Departments() error = %q, want token error", err)
	}
}
