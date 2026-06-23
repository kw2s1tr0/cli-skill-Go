package login

import (
	rootrepository "aiagentcliapp/repository"
	repositorylogin "aiagentcliapp/repository/login"
	loginrequest "aiagentcliapp/repository/login/request"
	"aiagentcliapp/service/login/input"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var actual loginrequest.Request
		if err := json.NewDecoder(request.Body).Decode(&actual); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if actual.Email != "user@example.com" || actual.Password != "password" || actual.TokenName != "agent-cli" {
			t.Fatalf("request = %#v", actual)
		}

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"access_token": "1|secret-token",
			"token_type": "Bearer",
			"expires_at": "2026-07-20T12:00:00+00:00"
		}`))
	}))
	defer server.Close()

	client := rootrepository.NewClient(server.URL, server.Client())
	service := NewService(repositorylogin.NewRepository(client))
	loginInput := input.Input{
		Email:     "user@example.com",
		Password:  "password",
		TokenName: "agent-cli",
	}

	if err := service.Login(context.Background(), loginInput); err != nil {
		t.Fatalf("Login() error = %v", err)
	}
}
