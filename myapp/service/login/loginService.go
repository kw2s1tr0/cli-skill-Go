package login

import (
	"aiagentcliapp/service/login/input"
	"aiagentcliapp/service/login/output"
	"context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (service *Service) Login(
	ctx context.Context,
	loginInput input.Input,
) (output.Output, error) {
	return output.Output{}, nil
}
