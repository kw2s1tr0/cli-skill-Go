package me

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/me/response"
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

func (repository *Repository) Me(ctx context.Context, accessToken string) (response.Response, error) {
	var user response.Response

	err := repository.client.DoJSON(
		ctx,
		http.MethodGet,
		"/api/auth/me",
		accessToken,
		nil,
		&user,
	)
	if err != nil {
		return response.Response{}, fmt.Errorf("me: %w", err)
	}

	return user, nil
}
