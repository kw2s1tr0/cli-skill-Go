package builder

import (
	"aiagentcliapp/service/department/input"
	"flag"
	"io"
	"strings"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) Build(args []string) (input.Input, error) {
	flags := flag.NewFlagSet("departments", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.Usage = func() {}

	orderBy := flags.String("order-by", "", "order by id, code, or name")
	orderDirection := flags.String("order-direction", "", "order direction asc or desc")
	if err := flags.Parse(args); err != nil {
		return input.Input{}, err
	}

	resolvedOrderBy := strings.TrimSpace(*orderBy)
	resolvedOrderDirection := strings.TrimSpace(*orderDirection)

	return input.Input{
		OrderBy:        resolvedOrderBy,
		OrderDirection: resolvedOrderDirection,
	}, nil
}
