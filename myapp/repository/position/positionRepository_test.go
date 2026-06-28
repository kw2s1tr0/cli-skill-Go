package position

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/position/output"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPositions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", request.Method)
		}
		if request.URL.Path != "/api/positions" {
			t.Errorf("path = %s, want /api/positions", request.URL.Path)
		}
		if got := request.URL.Query().Get("order_by"); got != "code" {
			t.Errorf("order_by = %q, want code", got)
		}
		if got := request.URL.Query().Get("order_direction"); got != "asc" {
			t.Errorf("order_direction = %q, want asc", got)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]output.Output{
			{ID: 1, Code: "MGR", Name: "Manager"},
			{ID: 2, Code: "ENG", Name: "Engineer"},
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	actual, err := NewRepository(client).Positions(context.Background(), "1|secret-token", SearchInput{
		OrderBy:        "code",
		OrderDirection: "asc",
	})
	if err != nil {
		t.Fatalf("Positions() error = %v", err)
	}
	if len(actual) != 2 || actual[0].ID != 1 || actual[0].Code != "MGR" || actual[0].Name != "Manager" || actual[1].ID != 2 || actual[1].Code != "ENG" || actual[1].Name != "Engineer" {
		t.Fatalf("Positions() = %#v", actual)
	}
}

func TestPositionsReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	_, err := NewRepository(client).Positions(context.Background(), "1|secret-token", SearchInput{})
	if err == nil {
		t.Fatal("Positions() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "positions: API returned 500 Internal Server Error: server error") {
		t.Fatalf("Positions() error = %q, want API error", err)
	}
}
