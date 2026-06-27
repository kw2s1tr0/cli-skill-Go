package builder

import (
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	args := []string{
		"--email", " user@example.com ",
		"--token-name", " agent-cli ",
	}

	actual, err := NewBuilder().Build(args, "secret")
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.Email != "user@example.com" || actual.Password != "secret" || actual.TokenName != "agent-cli" {
		t.Fatalf("Build() = %#v", actual)
	}
}

func TestBuildAllowsOmittedFlags(t *testing.T) {
	actual, err := NewBuilder().Build(nil, "secret")
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.Email != "" || actual.Password != "secret" || actual.TokenName != "" {
		t.Fatalf("Build() = %#v, want zero value", actual)
	}
}

func TestBuildRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "password flag",
			args:    []string{"--password", "password"},
			wantErr: "flag provided but not defined",
		},
		{
			name:    "password flag with value",
			args:    []string{"--password=password"},
			wantErr: "flag provided but not defined",
		},
		{
			name:    "unknown flag",
			args:    []string{"--unknown", "value"},
			wantErr: "flag provided but not defined",
		},
		{
			name:    "flag value is missing",
			args:    []string{"--email"},
			wantErr: "flag needs an argument",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewBuilder().Build(test.args, "secret")
			if err == nil {
				t.Fatal("Build() error = nil, want error")
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("Build() error = %q, want it to contain %q", err, test.wantErr)
			}
		})
	}
}
