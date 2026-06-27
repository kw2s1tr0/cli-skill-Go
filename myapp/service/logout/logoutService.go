package logout

import (
	repositorylogout "aiagentcliapp/repository/logout"
	"context"
	"fmt"
)

type Service struct {
	repository *repositorylogout.Repository
}

func NewService(repository *repositorylogout.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) Logout(ctx context.Context) error {
	err := service.repository.Logout(ctx)
	if err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	return nil
}
