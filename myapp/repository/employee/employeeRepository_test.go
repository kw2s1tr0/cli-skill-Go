package employee

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/employee/output"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmployees(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", request.Method)
		}
		if request.URL.Path != "/api/employees" {
			t.Errorf("path = %s, want /api/employees", request.URL.Path)
		}
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

	client := repository.NewClient(server.URL, server.Client())
	actual, err := NewRepository(client).Employees(context.Background(), "1|secret-token")
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

func TestEmployeesReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	_, err := NewRepository(client).Employees(context.Background(), "1|secret-token")
	if err == nil {
		t.Fatal("Employees() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "employees: API returned 500 Internal Server Error: server error") {
		t.Fatalf("Employees() error = %q, want API error", err)
	}
}
