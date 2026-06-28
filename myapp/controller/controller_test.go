package controller

import (
	"aiagentcliapp/repository"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type fakeTokenStore struct {
	accessToken string
	deleted     bool
}

func (store *fakeTokenStore) SaveAccessToken(ctx context.Context, accessToken string) error {
	store.accessToken = accessToken
	return nil
}

func (store *fakeTokenStore) GetAccessToken(ctx context.Context) (string, error) {
	return store.accessToken, nil
}

func (store *fakeTokenStore) DeleteAccessToken(ctx context.Context) error {
	store.deleted = true
	return nil
}

func TestRunWithoutCommandShowsUsage(t *testing.T) {
	client := repository.NewClient("http://example.com", http.DefaultClient)
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		nil,
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitUsage {
		t.Fatalf("Run() = %d, want %d", code, ExitUsage)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want empty", got)
	}
	for _, want := range []string{
		"login [--email EMAIL] [--token-name NAME]  Login and save access token",
		"me                                         Show current user",
		"logout                                     Logout and delete access token",
		"departments [--order-by id|code|name] [--order-direction asc|desc]",
		"employees [--keyword KEYWORD] [--department-id ID] [--position-id ID] [--employment-status active|leave|retired]",
		"positions [--order-by id|code|name] [--order-direction asc|desc]",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("stderr = %q, want %q", stderr.String(), want)
		}
	}
}

func TestHelpShowsCommandUsage(t *testing.T) {
	client := repository.NewClient("http://example.com", http.DefaultClient)
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"help"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitOK {
		t.Fatalf("Run() = %d, want %d", code, ExitOK)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
	for _, want := range []string{
		"login [--email EMAIL] [--token-name NAME]  Login and save access token",
		"me                                         Show current user",
		"logout                                     Logout and delete access token",
		"departments [--order-by id|code|name] [--order-direction asc|desc]",
		"employees [--keyword KEYWORD] [--department-id ID] [--position-id ID] [--employment-status active|leave|retired]",
		"positions [--order-by id|code|name] [--order-direction asc|desc]",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q, want %q", stdout.String(), want)
		}
	}
}

func TestLoginPromptsForPassword(t *testing.T) {
	tokenStore := &fakeTokenStore{}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	var actual struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		TokenName string `json:"token_name"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/auth/login" {
			t.Errorf("path = %s, want /api/auth/login", request.URL.Path)
		}
		if err := json.NewDecoder(request.Body).Decode(&actual); err != nil {
			t.Errorf("Decode() error = %v", err)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]any{
			"access_token": "token",
			"token_type":   "Bearer",
			"expires_at":   time.Now().UTC(),
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	passwordReaderCalled := false
	passwordReader := func(fd int) ([]byte, error) {
		passwordReaderCalled = true
		return []byte("secret"), nil
	}
	var stdout, stderr bytes.Buffer

	code := Run(
		[]string{"login", "--email", "user@example.com", "--token-name", "agent-cli"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitOK {
		t.Fatalf("Run() = %d, want %d; stderr = %q", code, ExitOK, stderr.String())
	}
	if !passwordReaderCalled {
		t.Fatal("ReadPassword() was not called")
	}
	if actual.Email != "user@example.com" || actual.Password != "secret" || actual.TokenName != "agent-cli" {
		t.Fatalf("request = %#v", actual)
	}
	if tokenStore.accessToken != "token" {
		t.Fatalf("saved access token = %q, want token", tokenStore.accessToken)
	}
	if got := stdout.String(); got != "login succeeded\n" {
		t.Fatalf("stdout = %q, want login succeeded", got)
	}
	if got := stderr.String(); got != "Password: \n" {
		t.Fatalf("stderr = %q, want password prompt", got)
	}
}

func TestLoginPasswordReaderFailure(t *testing.T) {
	client := repository.NewClient("http://example.com", http.DefaultClient)
	passwordReader := func(fd int) ([]byte, error) {
		return nil, errors.New("terminal unavailable")
	}
	var stdout, stderr bytes.Buffer

	code := Run(
		[]string{"login", "--email", "user@example.com"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitRuntime {
		t.Fatalf("Run() = %d, want %d", code, ExitRuntime)
	}
	if !strings.Contains(stderr.String(), "read password: terminal unavailable") {
		t.Fatalf("stderr = %q, want reader error", stderr.String())
	}
}

func TestLoginRejectsPasswordFlagAfterReadingPassword(t *testing.T) {
	client := repository.NewClient("http://example.com", http.DefaultClient)
	passwordReaderCalled := false
	passwordReader := func(fd int) ([]byte, error) {
		passwordReaderCalled = true
		return []byte("secret"), nil
	}
	var stdout, stderr bytes.Buffer

	code := Run(
		[]string{"login", "--email", "user@example.com", "--password", "secret"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitUsage {
		t.Fatalf("Run() = %d, want %d", code, ExitUsage)
	}
	if !passwordReaderCalled {
		t.Fatal("ReadPassword() was not called")
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("stderr = %q, want password flag error", stderr.String())
	}
}

func TestMe(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/auth/me" {
			t.Errorf("path = %s, want /api/auth/me", request.URL.Path)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]any{
			"id":    123,
			"name":  "User Name",
			"email": "user@example.com",
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"me"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitOK {
		t.Fatalf("Run() = %d, want %d; stderr = %q", code, ExitOK, stderr.String())
	}
	if got := stdout.String(); got != "ID: 123\nName: User Name\nEmail: user@example.com\n" {
		t.Fatalf("stdout = %q, want user output", got)
	}
}

func TestLogout(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/auth/logout" {
			t.Errorf("path = %s, want /api/auth/logout", request.URL.Path)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"logout"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitOK {
		t.Fatalf("Run() = %d, want %d; stderr = %q", code, ExitOK, stderr.String())
	}
	if !tokenStore.deleted {
		t.Fatal("DeleteAccessToken() was not called")
	}
	if got := stdout.String(); got != "logout succeeded\n" {
		t.Fatalf("stdout = %q, want logout succeeded", got)
	}
}

func TestDepartments(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/departments" {
			t.Errorf("path = %s, want /api/departments", request.URL.Path)
		}
		if got := request.URL.Query().Get("order_by"); got != "name" {
			t.Errorf("order_by = %q, want name", got)
		}
		if got := request.URL.Query().Get("order_direction"); got != "desc" {
			t.Errorf("order_direction = %q, want desc", got)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]map[string]any{
			{"id": 1, "code": "SALES", "name": "Sales"},
			{"id": 2, "code": "ENG", "name": "Engineering"},
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"departments", "--order-by", "name", "--order-direction", "desc"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitOK {
		t.Fatalf("Run() = %d, want %d; stderr = %q", code, ExitOK, stderr.String())
	}
	if got := stdout.String(); got != "ID: 1\nCode: SALES\nName: Sales\nID: 2\nCode: ENG\nName: Engineering\n" {
		t.Fatalf("stdout = %q, want departments output", got)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}

func TestDepartmentsReturnsAPIError(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/departments" {
			t.Errorf("path = %s, want /api/departments", request.URL.Path)
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"departments"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitRuntime {
		t.Fatalf("Run() = %d, want %d", code, ExitRuntime)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want empty", got)
	}
	if !strings.Contains(stderr.String(), "departments: departments: API returned 500 Internal Server Error: server error") {
		t.Fatalf("stderr = %q, want departments API error", stderr.String())
	}
}

func TestEmployees(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
		json.NewEncoder(writer).Encode([]map[string]any{
			{
				"id":                1,
				"employee_number":   "EMP-00001",
				"department_id":     10,
				"position_id":       20,
				"family_name":       "Yamada",
				"given_name":        "Taro",
				"family_name_kana":  "ヤマダ",
				"given_name_kana":   "タロウ",
				"email":             "yamada@example.com",
				"employment_status": "active",
				"department": map[string]any{
					"id":   10,
					"code": "DEV",
					"name": "Development",
				},
				"position": map[string]any{
					"id":   20,
					"code": "ENG",
					"name": "Engineer",
				},
			},
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"employees", "--keyword", "Yamada", "--department-id", "10", "--position-id", "20", "--employment-status", "active"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitOK {
		t.Fatalf("Run() = %d, want %d; stderr = %q", code, ExitOK, stderr.String())
	}
	if got := stdout.String(); got != "ID: 1\nEmployeeNumber: EMP-00001\nName: Yamada Taro\nNameKana: ヤマダ タロウ\nEmail: yamada@example.com\nEmploymentStatus: active\nDepartment: 10 DEV Development\nPosition: 20 ENG Engineer\n" {
		t.Fatalf("stdout = %q, want employees output", got)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}

func TestEmployeesReturnsAPIError(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/employees" {
			t.Errorf("path = %s, want /api/employees", request.URL.Path)
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"employees"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitRuntime {
		t.Fatalf("Run() = %d, want %d", code, ExitRuntime)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want empty", got)
	}
	if !strings.Contains(stderr.String(), "employees: employees: API returned 500 Internal Server Error: server error") {
		t.Fatalf("stderr = %q, want employees API error", stderr.String())
	}
}

func TestPositions(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/positions" {
			t.Errorf("path = %s, want /api/positions", request.URL.Path)
		}
		if got := request.URL.Query().Get("order_by"); got != "code" {
			t.Errorf("order_by = %q, want code", got)
		}
		if got := request.URL.Query().Get("order_direction"); got != "asc" {
			t.Errorf("order_direction = %q, want asc", got)
		}
		if authorization := request.Header.Get("Authorization"); authorization != "Bearer 1|secret-token" {
			t.Errorf("Authorization = %q, want bearer token", authorization)
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode([]map[string]any{
			{"id": 1, "code": "MGR", "name": "Manager"},
			{"id": 2, "code": "ENG", "name": "Engineer"},
		})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"positions", "--order-by", "code", "--order-direction", "asc"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitOK {
		t.Fatalf("Run() = %d, want %d; stderr = %q", code, ExitOK, stderr.String())
	}
	if got := stdout.String(); got != "ID: 1\nCode: MGR\nName: Manager\nID: 2\nCode: ENG\nName: Engineer\n" {
		t.Fatalf("stdout = %q, want positions output", got)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}

func TestPositionsReturnsAPIError(t *testing.T) {
	tokenStore := &fakeTokenStore{accessToken: "1|secret-token"}
	originalNewTokenStoreRepository := newTokenStoreRepository
	newTokenStoreRepository = func() accessTokenStore {
		return tokenStore
	}
	defer func() {
		newTokenStoreRepository = originalNewTokenStoreRepository
	}()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/positions" {
			t.Errorf("path = %s, want /api/positions", request.URL.Path)
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"message": "server error"})
	}))
	defer server.Close()

	client := repository.NewClient(server.URL, server.Client())
	var stdout, stderr bytes.Buffer
	passwordReader := func(fd int) ([]byte, error) {
		t.Fatal("ReadPassword() was called")
		return nil, nil
	}

	code := Run(
		[]string{"positions"},
		client,
		context.Background(),
		os.Stdin,
		&stdout,
		&stderr,
		passwordReader,
	)

	if code != ExitRuntime {
		t.Fatalf("Run() = %d, want %d", code, ExitRuntime)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want empty", got)
	}
	if !strings.Contains(stderr.String(), "positions: positions: API returned 500 Internal Server Error: server error") {
		t.Fatalf("stderr = %q, want positions API error", stderr.String())
	}
}
