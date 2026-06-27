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
	tokenStore tokenStore
}

type tokenStore interface {
	SaveAccessToken(context.Context, string) error
}

func NewService(repository *repositorylogin.Repository, tokenStore tokenStore) *Service {
	return &Service{
		repository: repository,
		tokenStore: tokenStore,
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

	loginResponse, err := service.repository.Login(ctx, loginRequest)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	// OS keystringに保存する
	if err := service.tokenStore.SaveAccessToken(ctx, loginResponse.AccessToken); err != nil {
		return fmt.Errorf("save access token: %w", err)
	}

	return nil
}
