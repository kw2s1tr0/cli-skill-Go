package employee

import (
	rootrepository "aiagentcliapp/repository"
	repositoryemployee "aiagentcliapp/repository/employee"
	employeeresponse "aiagentcliapp/repository/employee/response"
	"aiagentcliapp/service/employee/input"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeTokenStore struct {
	accessToken string
	err         error
}

func (store *fakeTokenStore) GetAccessToken(ctx context.Context) (string, error) {
	return store.accessToken, store.err
}

func TestEmployees(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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

	client := rootrepository.NewClient(server.URL, server.Client())
	service := NewService(repositoryemployee.NewRepository(client), &fakeTokenStore{accessToken: "1|secret-token"})

	actual, err := service.Employees(context.Background(), input.Input{})
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

func TestEmployeesReturnsTokenStoreError(t *testing.T) {
	service := NewService(
		repositoryemployee.NewRepository(rootrepository.NewClient("http://example.com", http.DefaultClient)),
		&fakeTokenStore{err: errors.New("missing token")},
	)

	_, err := service.Employees(context.Background(), input.Input{})
	if err == nil {
		t.Fatal("Employees() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "get access token: missing token") {
		t.Fatalf("Employees() error = %q, want token error", err)
	}
}
