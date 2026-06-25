package request

// RequestはログインAPIへ送るJSON本文を表す。
type Request struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	TokenName string `json:"token_name"`
}
