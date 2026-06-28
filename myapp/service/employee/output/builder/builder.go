package builder

import (
	employeeresponse "aiagentcliapp/repository/employee/response"
	"aiagentcliapp/service/employee/output"
)

type Builder struct{}

// NewBuilderはRepositoryのレスポンスからService出力を作るBuilderを返す。
func NewBuilder() *Builder {
	return &Builder{}
}

// Buildは社員APIレスポンスの一覧をService出力の一覧へ詰め替える。
func (builder *Builder) Build(responses []employeeresponse.Response) []output.Output {
	result := make([]output.Output, 0, len(responses))
	for _, response := range responses {
		result = append(result, builder.build(response))
	}

	return result
}

func (builder *Builder) build(response employeeresponse.Response) output.Output {
	return output.Output{
		ID:               response.ID,
		EmployeeNumber:   response.EmployeeNumber,
		FamilyName:       response.FamilyName,
		GivenName:        response.GivenName,
		FamilyNameKana:   response.FamilyNameKana,
		GivenNameKana:    response.GivenNameKana,
		Email:            response.Email,
		EmploymentStatus: response.EmploymentStatus,
		Department:       builder.buildRelated(response.Department),
		Position:         builder.buildRelated(response.Position),
	}
}

func (builder *Builder) buildRelated(response employeeresponse.RelatedResponse) output.RelatedOutput {
	return output.RelatedOutput{
		ID:   response.ID,
		Code: response.Code,
		Name: response.Name,
	}
}
