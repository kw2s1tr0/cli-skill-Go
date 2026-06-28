package builder

import (
	positionresponse "aiagentcliapp/repository/position/response"
	"testing"
)

func TestBuild(t *testing.T) {
	actual := NewBuilder().Build([]positionresponse.Response{
		{
			ID:        1,
			Code:      "MGR",
			Name:      "Manager",
			CreatedAt: "2026-01-01T00:00:00Z",
			UpdatedAt: "2026-01-02T00:00:00Z",
		},
		{ID: 2, Code: "ENG", Name: "Engineer"},
	})

	if len(actual) != 2 ||
		actual[0].ID != 1 ||
		actual[0].Code != "MGR" ||
		actual[0].Name != "Manager" ||
		actual[1].ID != 2 ||
		actual[1].Code != "ENG" ||
		actual[1].Name != "Engineer" {
		t.Fatalf("Build() = %#v", actual)
	}
}

func TestBuildWithSingleResponse(t *testing.T) {
	actual := NewBuilder().Build([]positionresponse.Response{{
		ID:        1,
		Code:      "MGR",
		Name:      "Manager",
		CreatedAt: "2026-01-01T00:00:00Z",
		UpdatedAt: "2026-01-02T00:00:00Z",
	}})

	if len(actual) != 1 ||
		actual[0].ID != 1 ||
		actual[0].Code != "MGR" ||
		actual[0].Name != "Manager" {
		t.Fatalf("Build() = %#v", actual)
	}
}
