package request

type Request struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	TokenName string `json:"token_name"`
}
