package builder

import (
	"testing"
	"time"
)

func TestBuild(t *testing.T) {
	expiresAt := time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC)
	actual := NewBuilder().Build("1|secret", "Bearer", expiresAt)

	if actual.AccessToken != "1|secret" || actual.TokenType != "Bearer" {
		t.Fatalf("Build() = %#v", actual)
	}
	if !actual.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("Build() expiry = %v, want %v", actual.ExpiresAt, expiresAt)
	}
}
