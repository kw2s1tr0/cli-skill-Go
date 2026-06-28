package department

import (
	"aiagentcliapp/repository"
	departmentrequest "aiagentcliapp/repository/department/request"
	departmentresponse "aiagentcliapp/repository/department/response"
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

func (repository *Repository) Departments(ctx context.Context, accessToken string, requestInput departmentrequest.Request) ([]departmentresponse.Response, error) {
	var departments []departmentresponse.Response
	query := url.Values{}
	if requestInput.OrderBy != "" {
		query.Set("order_by", requestInput.OrderBy)
	}
	if requestInput.OrderDirection != "" {
		query.Set("order_direction", requestInput.OrderDirection)
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
