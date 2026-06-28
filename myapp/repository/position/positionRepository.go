package position

import (
	"aiagentcliapp/repository"
	"aiagentcliapp/repository/position/output"
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

func (repository *Repository) Positions(ctx context.Context, accessToken string) ([]output.Output, error) {
	var positions []output.Output

	err := repository.client.DoJSON(
		ctx,
		http.MethodGet,
		"/api/positions",
		accessToken,
		nil,
		&positions,
	)
	if err != nil {
		return nil, fmt.Errorf("positions: %w", err)
	}

	return positions, nil
}
