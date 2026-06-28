package department

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/department/output"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDepartments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", request.Method)
		}
		if request.URL.Path != "/api/departments" {
			t.Errorf("path = %s, want /api/departments", request.URL.Path)
		}
		if got := request.URL.Query().Get("order_by"); got != "name" {
			t.Errorf("order_by = %q, want name", got)
		}
		if got := request.URL.Query().Get("order_direction"); got != "desc" {
			t.Errorf("order_direction = %q, want desc", got)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]output.Output{
			{ID: 1, Code: "SALES", Name: "Sales"},
			{ID: 2, Code: "ENG", Name: "Engineering"},
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	actual, err := NewRepository(client).Departments(context.Background(), "1|secret-token", SearchInput{
		OrderBy:        "name",
		OrderDirection: "desc",
	})
	if err != nil {
		t.Fatalf("Departments() error = %v", err)
	}
	if len(actual) != 2 || actual[0].ID != 1 || actual[0].Code != "SALES" || actual[0].Name != "Sales" || actual[1].ID != 2 || actual[1].Code != "ENG" || actual[1].Name != "Engineering" {
		t.Fatalf("Departments() = %#v", actual)
	}
}

func TestDepartmentsReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	_, err := NewRepository(client).Departments(context.Background(), "1|secret-token", SearchInput{})
	if err == nil {
		t.Fatal("Departments() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "departments: API returned 500 Internal Server Error: server error") {
		t.Fatalf("Departments() error = %q, want API error", err)
	}
}
