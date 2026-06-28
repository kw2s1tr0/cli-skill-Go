package input

// InputはCLIから受け取った部署検索条件をServiceへ渡すための型。
type Input struct {
	OrderBy        string
	OrderDirection string
}
