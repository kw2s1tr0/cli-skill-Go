package position

import (
	repositoryposition "aiagentcliapp/repository/position"
	"aiagentcliapp/repository/position/output"
	"context"
	"fmt"
)

type tokenStore interface {
	GetAccessToken(context.Context) (string, error)
}

type Service struct {
	repository *repositoryposition.Repository
	tokenStore tokenStore
}

func NewService(repository *repositoryposition.Repository, tokenStore tokenStore) *Service {
	return &Service{
		repository: repository,
		tokenStore: tokenStore,
	}
}

func (service *Service) Positions(ctx context.Context) ([]output.Output, error) {
	accessToken, err := service.tokenStore.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	positions, err := service.repository.Positions(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("positions: %w", err)
	}

	return positions, nil
}
