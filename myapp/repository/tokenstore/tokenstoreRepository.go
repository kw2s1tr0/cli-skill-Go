package tokenstore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	configDirName = ".aiagentcliapp"
	tokenFileName = "access-token"
)

var userHomeDir = os.UserHomeDir

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (repository *Repository) SaveAccessToken(ctx context.Context, accessToken string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("save access token: %w", err)
	}
	tokenPath, err := accessTokenPath()
	if err != nil {
		return fmt.Errorf("save access token: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		return fmt.Errorf("save access token: %w", err)
	}
	if err := os.WriteFile(tokenPath, []byte(accessToken), 0600); err != nil {
		return fmt.Errorf("save access token: %w", err)
	}
	return nil
}

func (repository *Repository) GetAccessToken(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}
	tokenPath, err := accessTokenPath()
	if err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}
	accessToken, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}
	return strings.TrimSpace(string(accessToken)), nil
}

func (repository *Repository) DeleteAccessToken(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("delete access token: %w", err)
	}
	tokenPath, err := accessTokenPath()
	if err != nil {
		return fmt.Errorf("delete access token: %w", err)
	}
	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete access token: %w", err)
	}
	return nil
}

func accessTokenPath() (string, error) {
	homeDir, err := userHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configDirName, tokenFileName), nil
}
