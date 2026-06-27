package tokenstore

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestSaveAccessToken(t *testing.T) {
	originalSetPassword := setPassword
	defer func() {
		setPassword = originalSetPassword
	}()

	var actualService, actualAccount, actualPassword string
	setPassword = func(service, account, password string) error {
		actualService = service
		actualAccount = account
		actualPassword = password
		return nil
	}

	if err := NewRepository().SaveAccessToken(context.Background(), "1|secret-token"); err != nil {
		t.Fatalf("SaveAccessToken() error = %v", err)
	}
	if actualService != serviceName || actualAccount != accountName || actualPassword != "1|secret-token" {
		t.Fatalf("saved credential = %q, %q, %q", actualService, actualAccount, actualPassword)
	}
}

func TestSaveAccessTokenError(t *testing.T) {
	originalSetPassword := setPassword
	defer func() {
		setPassword = originalSetPassword
	}()

	setPassword = func(service, account, password string) error {
		return errors.New("keyring unavailable")
	}

	err := NewRepository().SaveAccessToken(context.Background(), "1|secret-token")
	if err == nil {
		t.Fatal("SaveAccessToken() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "save access token: keyring unavailable") {
		t.Fatalf("SaveAccessToken() error = %q", err)
	}
}
