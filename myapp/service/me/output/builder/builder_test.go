package builder

import (
	meresponse "aiagentcliapp/repository/me/response"
	"testing"
)

func TestBuildList(t *testing.T) {
	actual := NewBuilder().BuildList([]meresponse.Response{
		{ID: 1, Name: "Taro Yamada", Email: "taro@example.com"},
		{ID: 2, Name: "Jiro Yamada", Email: "jiro@example.com"},
	})

	if len(actual) != 2 ||
		actual[0].ID != 1 ||
		actual[0].Name != "Taro Yamada" ||
		actual[0].Email != "taro@example.com" ||
		actual[1].ID != 2 ||
		actual[1].Name != "Jiro Yamada" ||
		actual[1].Email != "jiro@example.com" {
		t.Fatalf("BuildList() = %#v", actual)
	}
}

func TestBuildListWithSingleResponse(t *testing.T) {
	actual := NewBuilder().BuildList([]meresponse.Response{{
		ID:    1,
		Name:  "Taro Yamada",
		Email: "taro@example.com",
	}})

	if len(actual) != 1 ||
		actual[0].ID != 1 ||
		actual[0].Name != "Taro Yamada" ||
		actual[0].Email != "taro@example.com" {
		t.Fatalf("BuildList() = %#v", actual)
	}
}
