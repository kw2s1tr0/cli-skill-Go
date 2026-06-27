package tokenstore

import (
	"context"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "aiagentcliapp"
	accountName = "access-token"
)

var setPassword = keyring.Set

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (repository *Repository) SaveAccessToken(ctx context.Context, accessToken string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("save access token: %w", err)
	}
	if err := setPassword(serviceName, accountName, accessToken); err != nil {
		return fmt.Errorf("save access token: %w", err)
	}
	return nil
}
