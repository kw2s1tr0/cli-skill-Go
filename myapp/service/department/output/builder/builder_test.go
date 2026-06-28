package builder

import (
	departmentresponse "aiagentcliapp/repository/department/response"
	"testing"
)

func TestBuild(t *testing.T) {
	actual := NewBuilder().Build([]departmentresponse.Response{
		{
			ID:        1,
			Code:      "SALES",
			Name:      "Sales",
			CreatedAt: "2026-01-01T00:00:00Z",
			UpdatedAt: "2026-01-02T00:00:00Z",
		},
		{ID: 2, Code: "ENG", Name: "Engineering"},
	})

	if len(actual) != 2 ||
		actual[0].ID != 1 ||
		actual[0].Code != "SALES" ||
		actual[0].Name != "Sales" ||
		actual[1].ID != 2 ||
		actual[1].Code != "ENG" ||
		actual[1].Name != "Engineering" {
		t.Fatalf("Build() = %#v", actual)
	}
}

func TestBuildWithSingleResponse(t *testing.T) {
	actual := NewBuilder().Build([]departmentresponse.Response{{
		ID:        1,
		Code:      "SALES",
		Name:      "Sales",
		CreatedAt: "2026-01-01T00:00:00Z",
		UpdatedAt: "2026-01-02T00:00:00Z",
	}})

	if len(actual) != 1 ||
		actual[0].ID != 1 ||
		actual[0].Code != "SALES" ||
		actual[0].Name != "Sales" {
		t.Fatalf("Build() = %#v", actual)
	}
}
