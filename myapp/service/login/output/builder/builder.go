package builder

import (
	"aiagentcliapp/service/login/output"
	"time"
)

// Build creates the output returned by the login service.
func Build(
	accessToken string,
	tokenType string,
	expiresAt time.Time,
) output.Output {
	return output.Output{
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpiresAt:   expiresAt,
	}
}
