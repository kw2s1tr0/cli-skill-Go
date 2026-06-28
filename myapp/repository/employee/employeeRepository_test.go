package employee

import (
	"aiagentcliapp/repository"
	employeerequest "aiagentcliapp/repository/employee/request"
	employeeresponse "aiagentcliapp/repository/employee/response"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmployees(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", request.Method)
		}
		if request.URL.Path != "/api/employees" {
			t.Errorf("path = %s, want /api/employees", request.URL.Path)
		}
		if got := request.URL.Query().Get("keyword"); got != "Yamada" {
			t.Errorf("keyword = %q, want Yamada", got)
		}
		if got := request.URL.Query().Get("department_id"); got != "10" {
			t.Errorf("department_id = %q, want 10", got)
		}
		if got := request.URL.Query().Get("position_id"); got != "20" {
			t.Errorf("position_id = %q, want 20", got)
		}
		if got := request.URL.Query().Get("employment_status"); got != "active" {
			t.Errorf("employment_status = %q, want active", got)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]employeeresponse.Response{
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
				Department:       employeeresponse.RelatedResponse{ID: 10, Code: "DEV", Name: "Development"},
				Position:         employeeresponse.RelatedResponse{ID: 20, Code: "ENG", Name: "Engineer"},
			},
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	actual, err := NewRepository(client).Employees(context.Background(), "1|secret-token", employeerequest.Request{
		Keyword:          "Yamada",
		DepartmentID:     "10",
		PositionID:       "20",
		EmploymentStatus: "active",
	})
	if err != nil {
		t.Fatalf("Employees() error = %v", err)
	}
	if len(actual) != 1 ||
		actual[0].ID != 1 ||
		actual[0].EmployeeNumber != "EMP-00001" ||
		actual[0].FamilyName != "Yamada" ||
		actual[0].GivenName != "Taro" ||
		actual[0].Email != "yamada@example.com" ||
		actual[0].EmploymentStatus != "active" ||
		actual[0].Department.ID != 10 ||
		actual[0].Department.Code != "DEV" ||
		actual[0].Department.Name != "Development" ||
		actual[0].Position.ID != 20 ||
		actual[0].Position.Code != "ENG" ||
		actual[0].Position.Name != "Engineer" {
		t.Fatalf("Employees() = %#v", actual)
	}
}

func TestEmployeesReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	_, err := NewRepository(client).Employees(context.Background(), "1|secret-token", employeerequest.Request{})
	if err == nil {
		t.Fatal("Employees() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "employees: API returned 500 Internal Server Error: server error") {
		t.Fatalf("Employees() error = %q, want API error", err)
	}
}
