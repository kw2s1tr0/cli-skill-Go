package login

import (
	repositorylogin "aiagentcliapp/repository/login"
	requestbuilder "aiagentcliapp/repository/login/request/builder"
	"aiagentcliapp/service/login/input"
	"context"
	"fmt"
)

type Service struct {
	repository *repositorylogin.Repository
}

func NewService(repository *repositorylogin.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) Login(
	ctx context.Context,
	loginInput input.Input,
) error {

	loginRequest := requestbuilder.NewBuilder().Build(
		loginInput.Email,
		loginInput.Password,
		loginInput.TokenName,
	)

	_, err := service.repository.Login(ctx, loginRequest)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	return nil
}
