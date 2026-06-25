package controller

import (
	"aiagentcliapp/repository"
	"context"
	"fmt"
	"io"
)

const (
	// 終了コードは、CLIを呼び出したシェルやAIが結果を機械的に判断するための値
	ExitOK      = 0
	ExitRuntime = 1
	ExitUsage   = 2
)

func Run(
	args []string,
	client *repository.Client,
	ctx context.Context,
	stdout, stderr io.Writer,
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
	case "me":
	case "logout":
	case "departments":
	case "employees":
	case "positions":
	default:
		printUnknownCommand(stderr, cmd)
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
	login [--json]
	me [--json]
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

// 未対応のコマンド名はそのまま表示し、利用者が入力ミスを見つけやすくする。
func printUnknownCommand(writer io.Writer, cmd string) {
	fmt.Fprintf(writer, "error: unknown command %q\n\n", cmd)
}

// エラー表示を1か所にまとめ、stdoutとstderrの使い分けを崩さないようにする。
func printError(writer io.Writer, err error) {
	fmt.Fprintln(writer, "error:", err)
}
