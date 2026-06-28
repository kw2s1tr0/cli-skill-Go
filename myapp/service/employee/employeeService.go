package employee

import (
	repositoryemployee "aiagentcliapp/repository/employee"
	"aiagentcliapp/repository/employee/output"
	"context"
	"fmt"
)

type tokenStore interface {
	GetAccessToken(context.Context) (string, error)
}

type Service struct {
	repository *repositoryemployee.Repository
	tokenStore tokenStore
}

func NewService(repository *repositoryemployee.Repository, tokenStore tokenStore) *Service {
	return &Service{
		repository: repository,
		tokenStore: tokenStore,
	}
}

func (service *Service) Employees(ctx context.Context, input repositoryemployee.SearchInput) ([]output.Output, error) {
	accessToken, err := service.tokenStore.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	employees, err := service.repository.Employees(ctx, accessToken, input)
	if err != nil {
		return nil, fmt.Errorf("employees: %w", err)
	}

	return employees, nil
}
