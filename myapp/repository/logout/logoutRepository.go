package logout

import (
	"aiagentcliapp/repository"
	"context"
	"fmt"
	"net/http"
)

type tokenStore interface {
	GetAccessToken(context.Context) (string, error)
	DeleteAccessToken(context.Context) error
}

type Repository struct {
	client     *repository.Client
	tokenStore tokenStore
}

func NewRepository(client *repository.Client, tokenStore tokenStore) *Repository {
	return &Repository{
		client:     client,
		tokenStore: tokenStore,
	}
}

func (repository *Repository) Logout(ctx context.Context) error {
	accessToken, err := repository.tokenStore.GetAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("get access token: %w", err)
	}

	err = repository.client.DoJSON(
		ctx,
		http.MethodPost,
		"/api/auth/logout",
		accessToken,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("logout: %w", err)
	}
	if err := repository.tokenStore.DeleteAccessToken(ctx); err != nil {
		return fmt.Errorf("delete access token: %w", err)
	}

	return nil
}
