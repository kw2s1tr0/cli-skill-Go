package builder

import "testing"

func TestBuild(t *testing.T) {
	actual := NewBuilder().Build("user@example.com", "password", "agent-cli")

	if actual.Email != "user@example.com" || actual.Password != "password" || actual.TokenName != "agent-cli" {
		t.Fatalf("Build() = %#v", actual)
	}
}
