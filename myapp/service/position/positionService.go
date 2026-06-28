package position

import (
	repositoryposition "aiagentcliapp/repository/position"
	positionrequestbuilder "aiagentcliapp/repository/position/request/builder"
	"aiagentcliapp/service/position/input"
	"aiagentcliapp/service/position/output"
	positionoutputbuilder "aiagentcliapp/service/position/output/builder"
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

func (service *Service) Positions(ctx context.Context, searchInput input.Input) ([]output.Output, error) {
	accessToken, err := service.tokenStore.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	request := positionrequestbuilder.NewBuilder().Build(searchInput.OrderBy, searchInput.OrderDirection)
	positions, err := service.repository.Positions(ctx, accessToken, request)
	if err != nil {
		return nil, fmt.Errorf("positions: %w", err)
	}

	return positionoutputbuilder.NewBuilder().BuildList(positions), nil
}
