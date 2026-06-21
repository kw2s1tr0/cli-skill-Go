package builder

import (
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	args := []string{
		"--email", " user@example.com ",
		"--password", "password",
		"--token-name", " agent-cli ",
	}

	actual, err := NewBuilder().Build(args)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.Email != "user@example.com" || actual.Password != "password" || actual.TokenName != "agent-cli" {
		t.Fatalf("Build() = %#v", actual)
	}
}

func TestBuildAllowsOmittedFlags(t *testing.T) {
	actual, err := NewBuilder().Build(nil)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if actual.Email != "" || actual.Password != "" || actual.TokenName != "" {
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
			name:    "provided email is empty",
			args:    []string{"--email", ""},
			wantErr: "--email must not be empty",
		},
		{
			name:    "provided password is empty",
			args:    []string{"--password", ""},
			wantErr: "--password must not be empty",
		},
		{
			name:    "provided token name is empty",
			args:    []string{"--token-name", ""},
			wantErr: "--token-name must not be empty",
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
		{
			name: "positional argument",
			args: []string{
				"--email", "user@example.com",
				"--password", "password",
				"--token-name", "agent-cli",
				"extra",
			},
			wantErr: "does not accept positional arguments",
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
