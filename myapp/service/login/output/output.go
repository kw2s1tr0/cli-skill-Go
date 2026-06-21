package output

import "time"

// Output contains the token issued by a successful login.
type Output struct {
	AccessToken string
	TokenType   string
	ExpiresAt   time.Time
}
