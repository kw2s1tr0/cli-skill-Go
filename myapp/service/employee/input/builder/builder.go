package builder

import (
	"aiagentcliapp/service/employee/input"
	"flag"
	"io"
	"strings"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) Build(args []string) (input.Input, error) {
	flags := flag.NewFlagSet("employees", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.Usage = func() {}

	keyword := flags.String("keyword", "", "search keyword")
	departmentID := flags.String("department-id", "", "department id")
	positionID := flags.String("position-id", "", "position id")
	employmentStatus := flags.String("employment-status", "", "employment status active, leave, or retired")
	if err := flags.Parse(args); err != nil {
		return input.Input{}, err
	}

	resolvedKeyword := strings.TrimSpace(*keyword)
	resolvedDepartmentID := strings.TrimSpace(*departmentID)
	resolvedPositionID := strings.TrimSpace(*positionID)
	resolvedEmploymentStatus := strings.TrimSpace(*employmentStatus)

	return input.Input{
		Keyword:          resolvedKeyword,
		DepartmentID:     resolvedDepartmentID,
		PositionID:       resolvedPositionID,
		EmploymentStatus: resolvedEmploymentStatus,
	}, nil
}
