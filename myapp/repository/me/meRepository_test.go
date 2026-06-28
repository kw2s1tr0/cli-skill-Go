package me

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/me/response"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", request.Method)
		}
		if request.URL.Path != "/api/auth/me" {
			t.Errorf("path = %s, want /api/auth/me", request.URL.Path)
		}
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

	client := repository.NewClient(server.URL, server.Client())
	actual, err := NewRepository(client).Me(context.Background(), "1|secret-token")
	if err != nil {
		t.Fatalf("Me() error = %v", err)
	}
	if actual.ID != 123 || actual.Name != "User Name" || actual.Email != "user@example.com" {
		t.Fatalf("Me() = %#v", actual)
	}
}
