package builder

import "aiagentcliapp/repository/employee/request"

type Builder struct{}

// NewBuilderはServiceの値から社員検索APIリクエスト型を作るBuilderを返す。
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildはService層の社員検索入力を、Repository層のRequestへ詰め替える。
func (builder *Builder) Build(
	keyword string,
	departmentID string,
	positionID string,
	employmentStatus string,
) request.Request {
	return request.Request{
		Keyword:          keyword,
		DepartmentID:     departmentID,
		PositionID:       positionID,
		EmploymentStatus: employmentStatus,
	}
}
