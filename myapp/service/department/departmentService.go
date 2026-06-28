package department

import (
	repositorydepartment "aiagentcliapp/repository/department"
	departmentrequestbuilder "aiagentcliapp/repository/department/request/builder"
	"aiagentcliapp/service/department/input"
	"aiagentcliapp/service/department/output"
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

func (service *Service) Departments(ctx context.Context, searchInput input.Input) ([]output.Output, error) {
	accessToken, err := service.tokenStore.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	request := departmentrequestbuilder.NewBuilder().Build(searchInput.OrderBy, searchInput.OrderDirection)
	departments, err := service.repository.Departments(ctx, accessToken, request)
	if err != nil {
		return nil, fmt.Errorf("departments: %w", err)
	}

	result := make([]output.Output, 0, len(departments))
	for _, department := range departments {
		result = append(result, output.Output{
			ID:   department.ID,
			Code: department.Code,
			Name: department.Name,
		})
	}

	return result, nil
}
