package builder

import "aiagentcliapp/repository/department/request"

type Builder struct{}

// NewBuilderはServiceの値から部署検索APIリクエスト型を作るBuilderを返す。
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildはService層の部署検索入力を、Repository層のRequestへ詰め替える。
func (builder *Builder) Build(
	orderBy string,
	orderDirection string,
) request.Request {
	return request.Request{
		OrderBy:        orderBy,
		OrderDirection: orderDirection,
	}
}
