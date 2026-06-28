package output

type Output struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
