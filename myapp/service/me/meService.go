package me

import (
	repositoryme "aiagentcliapp/repository/me"
	meresponse "aiagentcliapp/repository/me/response"
	"aiagentcliapp/service/me/output"
	meoutputbuilder "aiagentcliapp/service/me/output/builder"
	"context"
	"fmt"
)

type tokenStore interface {
	GetAccessToken(context.Context) (string, error)
}

type Service struct {
	repository *repositoryme.Repository
	tokenStore tokenStore
}

func NewService(repository *repositoryme.Repository, tokenStore tokenStore) *Service {
	return &Service{
		repository: repository,
		tokenStore: tokenStore,
	}
}

func (service *Service) Me(ctx context.Context) (output.Output, error) {
	accessToken, err := service.tokenStore.GetAccessToken(ctx)
	if err != nil {
		return output.Output{}, fmt.Errorf("get access token: %w", err)
	}

	user, err := service.repository.Me(ctx, accessToken)
	if err != nil {
		return output.Output{}, fmt.Errorf("me: %w", err)
	}

	users := meoutputbuilder.NewBuilder().BuildList([]meresponse.Response{user})
	return users[0], nil
}
