package builder

import (
	positionresponse "aiagentcliapp/repository/position/response"
	"aiagentcliapp/service/position/output"
)

type Builder struct{}

// NewBuilderはRepositoryのレスポンスからService出力を作るBuilderを返す。
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildListは役職APIレスポンスの一覧をService出力の一覧へ詰め替える。
func (builder *Builder) BuildList(responses []positionresponse.Response) []output.Output {
	result := make([]output.Output, 0, len(responses))
	for _, response := range responses {
		result = append(result, builder.build(response))
	}

	return result
}

func (builder *Builder) build(response positionresponse.Response) output.Output {
	return output.Output{
		ID:   response.ID,
		Code: response.Code,
		Name: response.Name,
	}
}
