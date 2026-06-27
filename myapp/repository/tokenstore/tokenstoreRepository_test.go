package tokenstore

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withTempHome(t *testing.T) string {
	t.Helper()

	tempHome := t.TempDir()
	originalUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return tempHome, nil
	}
	t.Cleanup(func() {
		userHomeDir = originalUserHomeDir
	})

	return tempHome
}

func TestSaveAccessToken(t *testing.T) {
	tempHome := withTempHome(t)

	if err := NewRepository().SaveAccessToken(context.Background(), "1|secret-token"); err != nil {
		t.Fatalf("SaveAccessToken() error = %v", err)
	}

	tokenPath := filepath.Join(tempHome, configDirName, tokenFileName)
	actual, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(actual) != "1|secret-token" {
		t.Fatalf("saved token = %q, want 1|secret-token", actual)
	}

	dirInfo, err := os.Stat(filepath.Dir(tokenPath))
	if err != nil {
		t.Fatalf("Stat(config dir) error = %v", err)
	}
	if got := dirInfo.Mode().Perm(); got != 0700 {
		t.Fatalf("config dir permission = %o, want 700", got)
	}

	fileInfo, err := os.Stat(tokenPath)
	if err != nil {
		t.Fatalf("Stat(token file) error = %v", err)
	}
	if got := fileInfo.Mode().Perm(); got != 0600 {
		t.Fatalf("token file permission = %o, want 600", got)
	}
}

func TestGetAccessToken(t *testing.T) {
	tempHome := withTempHome(t)
	tokenPath := filepath.Join(tempHome, configDirName, tokenFileName)
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(tokenPath, []byte("1|secret-token"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	actual, err := NewRepository().GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken() error = %v", err)
	}
	if actual != "1|secret-token" {
		t.Fatalf("GetAccessToken() = %q, want 1|secret-token", actual)
	}
}

func TestGetAccessTokenTrimsWhitespace(t *testing.T) {
	tempHome := withTempHome(t)
	tokenPath := filepath.Join(tempHome, configDirName, tokenFileName)
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(tokenPath, []byte("  1|secret-token\n"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	actual, err := NewRepository().GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken() error = %v", err)
	}
	if actual != "1|secret-token" {
		t.Fatalf("GetAccessToken() = %q, want trimmed token", actual)
	}
}

func TestDeleteAccessToken(t *testing.T) {
	tempHome := withTempHome(t)
	tokenPath := filepath.Join(tempHome, configDirName, tokenFileName)
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(tokenPath, []byte("1|secret-token"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := NewRepository().DeleteAccessToken(context.Background()); err != nil {
		t.Fatalf("DeleteAccessToken() error = %v", err)
	}
	if _, err := os.Stat(tokenPath); !os.IsNotExist(err) {
		t.Fatalf("Stat(token file) error = %v, want not exist", err)
	}
}

func TestDeleteAccessTokenIgnoresMissingFile(t *testing.T) {
	withTempHome(t)

	if err := NewRepository().DeleteAccessToken(context.Background()); err != nil {
		t.Fatalf("DeleteAccessToken() error = %v", err)
	}
}

func TestSaveAccessTokenReturnsHomeDirError(t *testing.T) {
	originalUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return "", os.ErrPermission
	}
	t.Cleanup(func() {
		userHomeDir = originalUserHomeDir
	})

	err := NewRepository().SaveAccessToken(context.Background(), "1|secret-token")
	if err == nil {
		t.Fatal("SaveAccessToken() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "save access token") {
		t.Fatalf("SaveAccessToken() error = %q", err)
	}
}
