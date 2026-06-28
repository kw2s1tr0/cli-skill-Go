package position

import (
	"aiagentcliapp/repository"
	positionrequest "aiagentcliapp/repository/position/request"
	positionresponse "aiagentcliapp/repository/position/response"
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

func (repository *Repository) Positions(ctx context.Context, accessToken string, requestInput positionrequest.Request) ([]positionresponse.Response, error) {
	var positions []positionresponse.Response
	query := url.Values{}
	if requestInput.OrderBy != "" {
		query.Set("order_by", requestInput.OrderBy)
	}
	if requestInput.OrderDirection != "" {
		query.Set("order_direction", requestInput.OrderDirection)
	}
	path := "/api/positions"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}

	err := repository.client.DoJSON(
		ctx,
		http.MethodGet,
		path,
		accessToken,
		nil,
		&positions,
	)
	if err != nil {
		return nil, fmt.Errorf("positions: %w", err)
	}

	return positions, nil
}
