package controller

import (
	"aiagentcliapp/repository"
	repositorylogin "aiagentcliapp/repository/login"
	repositorylogout "aiagentcliapp/repository/logout"
	repositoryme "aiagentcliapp/repository/me"
	repositorytokenstore "aiagentcliapp/repository/tokenstore"
	servicelogin "aiagentcliapp/service/login"
	inputbuilder "aiagentcliapp/service/login/input/builder"
	servicelogout "aiagentcliapp/service/logout"
	serviceme "aiagentcliapp/service/me"
	"context"
	"fmt"
	"io"
	"os"
)

const (
	// 終了コードは、CLIを呼び出したシェルやAIが結果を機械的に判断するための値
	ExitOK      = 0
	ExitRuntime = 1
	ExitUsage   = 2
)

// io.Writerのような型が無いため独自の関数の型をつくる
type PasswordReader func(fd int) ([]byte, error)

type accessTokenStore interface {
	SaveAccessToken(context.Context, string) error
	GetAccessToken(context.Context) (string, error)
	DeleteAccessToken(context.Context) error
}

var newTokenStoreRepository = func() accessTokenStore {
	return repositorytokenstore.NewRepository()
}

func Run(
	args []string,
	client *repository.Client,
	ctx context.Context,
	stdin *os.File,
	stdout, stderr io.Writer,
	passwordReader PasswordReader,
) int {
	// サブコマンドが無い場合は使い方の誤りとして扱い、使用方法をstderrへ出す。
	if len(args) == 0 {
		printUsage(stderr)
		return ExitUsage
	}

	// 先頭の引数をサブコマンド名として扱い、残りは各コマンドの引数にする。
	cmd := args[0]

	// 入力ミスと実行中の失敗を分けると、呼び出し側が終了コードで原因を判断できる。
	var runtimeErr error
	var validationErr error

	// Controllerはコマンド名を判定し、実際の業務処理はServiceへ委譲する。
	switch cmd {

	case "help", "-h", "--help":
		printHelp(stdout)
		return ExitOK

	case "login":
		// login入力を非表示で行う
		fmt.Fprint(stderr, "Password: ")
		password, err := passwordReader(int(stdin.Fd()))
		// 改行
		fmt.Fprintln(stderr)
		if err != nil {
			runtimeErr = fmt.Errorf("read password: %w", err)
			break
		}

		loginInput, err := inputbuilder.NewBuilder().Build(args[1:], string(password))
		if err != nil {
			validationErr = err
			break
		}

		// テスト容易性のためDIする
		loginRepository := repositorylogin.NewRepository(client)
		tokenStoreRepository := newTokenStoreRepository()
		loginService := servicelogin.NewService(loginRepository, tokenStoreRepository)

		runtimeErr = loginService.Login(ctx, loginInput)
		if runtimeErr == nil {
			fmt.Fprintln(stdout, "login succeeded")
		}

	case "me":
		meRepository := repositoryme.NewRepository(client)
		tokenStoreRepository := newTokenStoreRepository()
		meService := serviceme.NewService(meRepository, tokenStoreRepository)

		user, err := meService.Me(ctx)
		if err != nil {
			runtimeErr = err
			break
		}
		fmt.Fprintf(stdout, "ID: %d\nName: %s\nEmail: %s\n", user.ID, user.Name, user.Email)

	case "logout":
		tokenStoreRepository := newTokenStoreRepository()
		logoutRepository := repositorylogout.NewRepository(client, tokenStoreRepository)
		logoutService := servicelogout.NewService(logoutRepository)

		runtimeErr = logoutService.Logout(ctx)
		if runtimeErr == nil {
			fmt.Fprintln(stdout, "logout succeeded")
		}

	case "departments":

	case "employees":

	case "positions":

	default:
		// 想定外commandの場合はエラー出力
		fmt.Fprintf(stderr, "error: unknown command %q\n\n", cmd)
		printUsage(stderr)
		return ExitUsage
	}

	// 入力エラーは使用方法と一緒に表示し、実行時エラーとは別の終了コードにする。
	if validationErr != nil {
		printError(stderr, validationErr)
		printUsage(stderr)
		return ExitUsage
	}

	if runtimeErr != nil {
		printError(stderr, runtimeErr)
		return ExitRuntime
	}

	return ExitOK
}

// printUsageは、入力が足りない・間違っている場合に最低限の使い方を表示する。
func printUsage(writer io.Writer) {
	fmt.Fprintln(writer, `Usage:
	login [--email EMAIL] [--token-name NAME]
	me
	logout
	departments
	employees
	positions
	help`)
}

// printHelpは、明示的にhelpが指定されたときのコマンド一覧を表示する。
func printHelp(writer io.Writer) {
	fmt.Fprintln(writer, `command list:
	login
	me
	logout
	departments
	employees
	positions`)
}

// エラー表示を1か所にまとめ、stdoutとstderrの使い分けを崩さないようにする。
func printError(writer io.Writer, err error) {
	fmt.Fprintln(writer, "error:", err)
}
