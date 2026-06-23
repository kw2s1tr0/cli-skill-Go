package builder

import "aiagentcliapp/repository/login/request"

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) Build(
	email string,
	password string,
	tokenName string,
) request.Request {
	return request.Request{
		Email:     email,
		Password:  password,
		TokenName: tokenName,
	}
}
