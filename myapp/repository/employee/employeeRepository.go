package employee

import (
	"aiagentcliapp/repository"
	employeerequest "aiagentcliapp/repository/employee/request"
	employeeresponse "aiagentcliapp/repository/employee/response"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Repository struct {
	client *repository.Client
}

func NewRepository(client *repository.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repository *Repository) Employees(ctx context.Context, accessToken string, requestInput employeerequest.Request) ([]employeeresponse.Response, error) {
	var employees []employeeresponse.Response
	query := url.Values{}
	if requestInput.Keyword != "" {
		query.Set("keyword", requestInput.Keyword)
	}
	if requestInput.DepartmentID != "" {
		query.Set("department_id", requestInput.DepartmentID)
	}
	if requestInput.PositionID != "" {
		query.Set("position_id", requestInput.PositionID)
	}
	if requestInput.EmploymentStatus != "" {
		query.Set("employment_status", requestInput.EmploymentStatus)
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
