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
	if len(args) == 0 {
		printUsage(stderr)
		return ExitUsage
	}

	cmd := args[0]
	var runtimeErr error
	var validationErr error

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

func printHelp(writer io.Writer) {
	fmt.Fprintln(writer, `command list:
	login
	me
	logout
	departments
	employees
	positions`)
}

func printUnknownCommand(writer io.Writer, cmd string) {
	fmt.Fprintf(writer, "error: unknown command %q\n\n", cmd)
}

func printError(writer io.Writer, err error) {
	fmt.Fprintln(writer, "error:", err)
}
