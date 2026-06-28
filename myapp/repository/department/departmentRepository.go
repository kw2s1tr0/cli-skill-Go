package department

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/department/output"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Repository struct {
	client *repository.Client
}

type SearchInput struct {
	OrderBy        string
	OrderDirection string
}

func NewRepository(client *repository.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repository *Repository) Departments(ctx context.Context, accessToken string, input SearchInput) ([]output.Output, error) {
	var departments []output.Output
	query := url.Values{}
	if input.OrderBy != "" {
		query.Set("order_by", input.OrderBy)
	}
	if input.OrderDirection != "" {
		query.Set("order_direction", input.OrderDirection)
	}
	path := "/api/departments"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}

	err := repository.client.DoJSON(
		ctx,
		http.MethodGet,
		path,
		accessToken,
		nil,
		&departments,
	)
	if err != nil {
		return nil, fmt.Errorf("departments: %w", err)
	}

	return departments, nil
}
