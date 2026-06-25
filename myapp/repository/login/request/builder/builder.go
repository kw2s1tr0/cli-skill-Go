package builder

import "aiagentcliapp/repository/login/request"

type Builder struct{}

// NewBuilderはServiceの値からAPIリクエスト型を作るBuilderを返す。
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildはService層のログイン入力を、APIが要求するJSON用のRequestへ詰め替える。
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
