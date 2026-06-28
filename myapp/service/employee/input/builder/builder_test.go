package builder

import (
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	actual, err := NewBuilder().Build([]string{
		"--keyword", " Yamada ",
		"--department-id", " 10 ",
		"--position-id", " 20 ",
		"--employment-status", " active ",
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.Keyword != "Yamada" ||
		actual.DepartmentID != "10" ||
		actual.PositionID != "20" ||
		actual.EmploymentStatus != "active" {
		t.Fatalf("Build() = %#v", actual)
	}
}

func TestBuildAllowsOmittedFlags(t *testing.T) {
	actual, err := NewBuilder().Build(nil)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.Keyword != "" || actual.DepartmentID != "" || actual.PositionID != "" || actual.EmploymentStatus != "" {
		t.Fatalf("Build() = %#v, want zero value", actual)
	}
}

func TestBuildLeavesValidationToAPI(t *testing.T) {
	actual, err := NewBuilder().Build([]string{
		"--department-id", "abc",
		"--position-id", "0",
		"--employment-status", "unknown",
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.DepartmentID != "abc" || actual.PositionID != "0" || actual.EmploymentStatus != "unknown" {
		t.Fatalf("Build() = %#v", actual)
	}
}

func TestBuildRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "unknown flag",
			args:    []string{"--unknown", "value"},
			wantErr: "flag provided but not defined",
		},
		{
			name:    "flag value is missing",
			args:    []string{"--keyword"},
			wantErr: "flag needs an argument",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewBuilder().Build(test.args)
			if err == nil {
				t.Fatal("Build() error = nil, want error")
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("Build() error = %q, want it to contain %q", err, test.wantErr)
			}
		})
	}
}
