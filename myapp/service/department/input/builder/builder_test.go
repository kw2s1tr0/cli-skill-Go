package builder

import (
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	actual, err := NewBuilder().Build([]string{
		"--order-by", " name ",
		"--order-direction", " desc ",
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.OrderBy != "name" || actual.OrderDirection != "desc" {
		t.Fatalf("Build() = %#v", actual)
	}
}

func TestBuildAllowsOmittedFlags(t *testing.T) {
	actual, err := NewBuilder().Build(nil)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.OrderBy != "" || actual.OrderDirection != "" {
		t.Fatalf("Build() = %#v, want zero value", actual)
	}
}

func TestBuildLeavesValidationToAPI(t *testing.T) {
	actual, err := NewBuilder().Build([]string{
		"--order-by", "created_at",
		"--order-direction", "sideways",
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.OrderBy != "created_at" || actual.OrderDirection != "sideways" {
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
			args:    []string{"--order-by"},
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
