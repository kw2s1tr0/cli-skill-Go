package response

import "time"

// ResponseはログインAPIから返るトークン情報を表す。
type Response struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}
