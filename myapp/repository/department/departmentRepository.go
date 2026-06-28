package department

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/department/output"
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

func (repository *Repository) Departments(ctx context.Context, accessToken string) ([]output.Output, error) {
	var departments []output.Output

	err := repository.client.DoJSON(
		ctx,
		http.MethodGet,
		"/api/departments",
		accessToken,
		nil,
		&departments,
	)
	if err != nil {
		return nil, fmt.Errorf("departments: %w", err)
	}

	return departments, nil
}
