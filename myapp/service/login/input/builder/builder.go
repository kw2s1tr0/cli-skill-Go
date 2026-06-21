package builder

import (
	"aiagentcliapp/service/login/input"
	"flag"
	"io"
	"strings"
)

// command以下の引数を渡してserviceの引数に変換する
func Build(args []string) (input.Input, error) {
	// Go標準のCLIの引数解析ツール
	// 命名とエラー時に処理を停めないようにしている
	flags := flag.NewFlagSet("login", flag.ContinueOnError)
	// errorを画面に出さないようにする
	flags.SetOutput(io.Discard)
	// Usageを出さないようにする（これをしない場合--helpで自動で出してくれる）
	flags.Usage = func() {}

	// 引数を定義
	email := flags.String("email", "", "login email address")
	password := flags.String("password", "", "login password")
	tokenName := flags.String("token-name", "", "API token name")

	// 引数を検証
	if err := flags.Parse(args); err != nil {
		return input.Input{}, err
	}

	resolvedEmail := strings.TrimSpace(*email)
	resolvedTokenName := strings.TrimSpace(*tokenName)

	return input.Input{
		Email:     resolvedEmail,
		Password:  *password,
		TokenName: resolvedTokenName,
	}, nil
}
