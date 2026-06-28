package employee

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/employee/output"
	"context"
	"fmt"
	"net/http"
)

type Repository struct {
	client *repository.Client
}

func NewRepository(client *repository.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repository *Repository) Employees(ctx context.Context, accessToken string) ([]output.Output, error) {
	var employees []output.Output

	err := repository.client.DoJSON(
		ctx,
		http.MethodGet,
		"/api/employees",
		accessToken,
		nil,
		&employees,
	)
	if err != nil {
		return nil, fmt.Errorf("employees: %w", err)
	}

	return employees, nil
}
