package login

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/login/request"
	"aiagentcliapp/repository/login/response"
	"context"
	"fmt"
	"net/http"
)

type Repository struct {
	client *repository.Client
}

func NewRepository(client *repository.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repository *Repository) Login(
	ctx context.Context,
	loginRequest request.Request,
) (response.Response, error) {
	var loginResponse response.Response

	method := http.MethodPost
	path := "/api/auth/login"
	body := loginRequest
	result := &loginResponse

	err := repository.client.DoJSON(ctx, method, path, body, result)
	if err != nil {
		return response.Response{}, fmt.Errorf("login: %w", err)
	}

	return loginResponse, nil
}
