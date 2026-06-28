package employee

import (
	repositoryemployee "aiagentcliapp/repository/employee"
	employeerequestbuilder "aiagentcliapp/repository/employee/request/builder"
	"aiagentcliapp/service/employee/input"
	"aiagentcliapp/service/employee/output"
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

func (service *Service) Employees(ctx context.Context, searchInput input.Input) ([]output.Output, error) {
	accessToken, err := service.tokenStore.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	request := employeerequestbuilder.NewBuilder().Build(
		searchInput.Keyword,
		searchInput.DepartmentID,
		searchInput.PositionID,
		searchInput.EmploymentStatus,
	)
	employees, err := service.repository.Employees(ctx, accessToken, request)
	if err != nil {
		return nil, fmt.Errorf("employees: %w", err)
	}

	result := make([]output.Output, 0, len(employees))
	for _, employee := range employees {
		result = append(result, output.Output{
			ID:               employee.ID,
			EmployeeNumber:   employee.EmployeeNumber,
			FamilyName:       employee.FamilyName,
			GivenName:        employee.GivenName,
			FamilyNameKana:   employee.FamilyNameKana,
			GivenNameKana:    employee.GivenNameKana,
			Email:            employee.Email,
			EmploymentStatus: employee.EmploymentStatus,
			Department: output.RelatedOutput{
				ID:   employee.Department.ID,
				Code: employee.Department.Code,
				Name: employee.Department.Name,
			},
			Position: output.RelatedOutput{
				ID:   employee.Position.ID,
				Code: employee.Position.Code,
				Name: employee.Position.Name,
			},
		})
	}

	return result, nil
}
