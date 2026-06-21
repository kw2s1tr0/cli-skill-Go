package login

import (
	"aiagentcliapp/service/login/input"
	"context"
	"errors"
	"testing"
	"time"
)

type fakeTokenIssuer struct {
	email       string
	password    string
	tokenName   string
	accessToken string
	tokenType   string
	expiresAt   time.Time
	err         error
}

func (fake *fakeTokenIssuer) Issue(
	_ context.Context,
	email string,
	password string,
	tokenName string,
) (string, string, time.Time, error) {
	fake.email = email
	fake.password = password
	fake.tokenName = tokenName

	return fake.accessToken, fake.tokenType, fake.expiresAt, fake.err
}

func TestLogin(t *testing.T) {
	expiresAt := time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC)
	issuer := &fakeTokenIssuer{
		accessToken: "1|secret-token",
		tokenType:   "Bearer",
		expiresAt:   expiresAt,
	}
	service := NewService(issuer)
	loginInput := input.Input{
		Email:     "user@example.com",
		Password:  "password",
		TokenName: "agent-cli",
	}

	actual, err := service.Login(context.Background(), loginInput)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if actual.AccessToken != issuer.accessToken || actual.TokenType != issuer.tokenType {
		t.Fatalf("Login() = %#v", actual)
	}
	if !actual.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("Login() expiry = %v, want %v", actual.ExpiresAt, expiresAt)
	}
	if issuer.email != loginInput.Email || issuer.password != loginInput.Password || issuer.tokenName != loginInput.TokenName {
		t.Fatalf("Issue() input = %q, %q, %q", issuer.email, issuer.password, issuer.tokenName)
	}
}

func TestLoginWrapsIssuerError(t *testing.T) {
	issuerError := errors.New("request failed")
	service := NewService(&fakeTokenIssuer{err: issuerError})

	_, err := service.Login(context.Background(), input.Input{})
	if !errors.Is(err, issuerError) {
		t.Fatalf("Login() error = %v, want wrapped issuer error", err)
	}
}
