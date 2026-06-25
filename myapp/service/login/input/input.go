package input

// InputはCLIから受け取ったログイン情報をServiceへ渡すための型。
type Input struct {
	Email     string
	Password  string
	TokenName string
}
