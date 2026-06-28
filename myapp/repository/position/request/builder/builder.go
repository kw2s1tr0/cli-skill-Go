package builder

import "aiagentcliapp/repository/position/request"

type Builder struct{}

// NewBuilderはServiceの値から役職検索APIリクエスト型を作るBuilderを返す。
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildはService層の役職検索入力を、Repository層のRequestへ詰め替える。
func (builder *Builder) Build(
	orderBy string,
	orderDirection string,
) request.Request {
	return request.Request{
		OrderBy:        orderBy,
		OrderDirection: orderDirection,
	}
}
