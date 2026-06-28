package employee

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/employee/output"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Repository struct {
	client *repository.Client
}

type SearchInput struct {
	Keyword          string
	DepartmentID     string
	PositionID       string
	EmploymentStatus string
}

func NewRepository(client *repository.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repository *Repository) Employees(ctx context.Context, accessToken string, input SearchInput) ([]output.Output, error) {
	var employees []output.Output
	query := url.Values{}
	if input.Keyword != "" {
		query.Set("keyword", input.Keyword)
	}
	if input.DepartmentID != "" {
		query.Set("department_id", input.DepartmentID)
	}
	if input.PositionID != "" {
		query.Set("position_id", input.PositionID)
	}
	if input.EmploymentStatus != "" {
		query.Set("employment_status", input.EmploymentStatus)
	}
	path := "/api/employees"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}

	err := repository.client.DoJSON(
		ctx,
		http.MethodGet,
		path,
		accessToken,
		nil,
		&employees,
	)
	if err != nil {
		return nil, fmt.Errorf("employees: %w", err)
	}

	return employees, nil
}
