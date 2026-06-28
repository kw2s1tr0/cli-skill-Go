package builder

import (
	employeeresponse "aiagentcliapp/repository/employee/response"
	"testing"
)

func TestBuildList(t *testing.T) {
	actual := NewBuilder().BuildList([]employeeresponse.Response{
		{
			ID:               1,
			EmployeeNumber:   "EMP-00001",
			DepartmentID:     10,
			PositionID:       20,
			FamilyName:       "Yamada",
			GivenName:        "Taro",
			FamilyNameKana:   "ヤマダ",
			GivenNameKana:    "タロウ",
			Email:            "yamada@example.com",
			EmploymentStatus: "active",
			Department:       employeeresponse.RelatedResponse{ID: 10, Code: "SALES", Name: "Sales"},
			Position:         employeeresponse.RelatedResponse{ID: 20, Code: "MGR", Name: "Manager"},
		},
		{ID: 2, EmployeeNumber: "EMP-00002"},
	})

	if len(actual) != 2 ||
		actual[0].ID != 1 ||
		actual[0].EmployeeNumber != "EMP-00001" ||
		actual[0].FamilyName != "Yamada" ||
		actual[0].GivenName != "Taro" ||
		actual[0].FamilyNameKana != "ヤマダ" ||
		actual[0].GivenNameKana != "タロウ" ||
		actual[0].Email != "yamada@example.com" ||
		actual[0].EmploymentStatus != "active" ||
		actual[0].Department.ID != 10 ||
		actual[0].Department.Code != "SALES" ||
		actual[0].Department.Name != "Sales" ||
		actual[0].Position.ID != 20 ||
		actual[0].Position.Code != "MGR" ||
		actual[0].Position.Name != "Manager" ||
		actual[1].ID != 2 ||
		actual[1].EmployeeNumber != "EMP-00002" {
		t.Fatalf("BuildList() = %#v", actual)
	}
}

func TestBuildListWithSingleResponse(t *testing.T) {
	actual := NewBuilder().BuildList([]employeeresponse.Response{{
		ID:             1,
		EmployeeNumber: "EMP-00001",
	}})

	if len(actual) != 1 ||
		actual[0].ID != 1 ||
		actual[0].EmployeeNumber != "EMP-00001" {
		t.Fatalf("BuildList() = %#v", actual)
	}
}
