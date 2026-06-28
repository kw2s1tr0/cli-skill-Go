package builder

import (
	meresponse "aiagentcliapp/repository/me/response"
	"aiagentcliapp/service/me/output"
)

type Builder struct{}

// NewBuilderはRepositoryのレスポンスからService出力を作るBuilderを返す。
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildはユーザーAPIレスポンスの一覧をService出力の一覧へ詰め替える。
func (builder *Builder) Build(responses []meresponse.Response) []output.Output {
	result := make([]output.Output, 0, len(responses))
	for _, response := range responses {
		result = append(result, builder.build(response))
	}

	return result
}

func (builder *Builder) build(response meresponse.Response) output.Output {
	return output.Output{
		ID:    response.ID,
		Name:  response.Name,
		Email: response.Email,
	}
}
