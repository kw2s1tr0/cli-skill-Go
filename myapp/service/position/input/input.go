package input

// InputはCLIから受け取った役職検索条件をServiceへ渡すための型。
type Input struct {
	OrderBy        string
	OrderDirection string
}
