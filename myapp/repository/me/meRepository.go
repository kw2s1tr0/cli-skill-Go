package me

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/me/output"
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

func (repository *Repository) Me(ctx context.Context, accessToken string) (output.Output, error) {
	var user output.Output

	err := repository.client.DoJSON(
		ctx,
		http.MethodGet,
		"/api/auth/me",
		accessToken,
		nil,
		&user,
	)
	if err != nil {
		return output.Output{}, fmt.Errorf("me: %w", err)
	}

	return user, nil
}
