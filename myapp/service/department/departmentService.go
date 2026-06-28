package department

import (
	repositorydepartment "aiagentcliapp/repository/department"
	"aiagentcliapp/repository/department/output"
	"context"
	"fmt"
)

type tokenStore interface {
	GetAccessToken(context.Context) (string, error)
}

type Service struct {
	repository *repositorydepartment.Repository
	tokenStore tokenStore
}

func NewService(repository *repositorydepartment.Repository, tokenStore tokenStore) *Service {
	return &Service{
		repository: repository,
		tokenStore: tokenStore,
	}
}

func (service *Service) Departments(ctx context.Context, input repositorydepartment.SearchInput) ([]output.Output, error) {
	accessToken, err := service.tokenStore.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	departments, err := service.repository.Departments(ctx, accessToken, input)
	if err != nil {
		return nil, fmt.Errorf("departments: %w", err)
	}

	return departments, nil
}
