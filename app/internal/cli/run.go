package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"company-cli/internal/api"
)

const (
	// 終了コードは、CLIを呼び出したシェルやAIが結果を機械的に判断するための値です。
	ExitOK      = 0
	ExitRuntime = 1
	ExitUsage   = 2
)

// postsAPIはCLIがAPI層に求める操作だけを定義します。
// 具体的なClientではなくinterfaceに依存することで、テストではFakeへ交換できます。
type postsAPI interface {
	ListPosts(context.Context, int) ([]api.Post, error)
	GetPost(context.Context, int) (api.Post, error)
	CreatePost(context.Context, api.CreatePostInput) (api.Post, error)
}

// RunはCLIを実行し、プロセスの終了コードを返します。
// 引数と入出力をosパッケージから切り離して受け取るため、テストしやすい形です。
func Run(ctx context.Context, args []string, stdout, stderr io.Writer, client postsAPI) int {
	if len(args) == 0 {
		printUsage(stderr)
		return ExitUsage
	}

	// 先頭の引数をサブコマンド名として振り分けます。
	switch args[0] {
	case "help", "-h", "--help":
		printUsage(stdout)
		return ExitOK
	case "list":
		return runList(ctx, args[1:], stdout, stderr, client)
	case "get":
		return runGet(ctx, args[1:], stdout, stderr, client)
	case "create":
		return runCreate(ctx, args[1:], stdout, stderr, client)
	default:
		fmt.Fprintf(stderr, "error: unknown command %q\n\n", args[0])
		printUsage(stderr)
		return ExitUsage
	}
}

func runList(ctx context.Context, args []string, stdout, stderr io.Writer, client postsAPI) int {
	args, outputJSON := extractJSONFlag(args)
	flags := newFlagSet("list", stderr)
	// flag.Intの戻り値は*int（ポインタ）です。Parse後に*limitで実際の値を読みます。
	limit := flags.Int("limit", 10, "maximum number of posts (1-100)")
	if err := flags.Parse(args); err != nil {
		return flagErrorCode(err)
	}
	if flags.NArg() != 0 {
		return usageError(stderr, "list does not accept positional arguments")
	}
	if *limit < 1 || *limit > 100 {
		return usageError(stderr, "--limit must be between 1 and 100")
	}

	posts, err := client.ListPosts(ctx, *limit)
	if err != nil {
		return runtimeError(stderr, err)
	}
	if outputJSON {
		return writeJSON(stdout, stderr, posts)
	}

	// tabwriterは\t区切りの文字列を、人が読みやすい列へ揃えてくれます。
	writer := tabwriter.NewWriter(stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tUSER\tTITLE")
	for _, post := range posts {
		fmt.Fprintf(writer, "%d\t%d\t%s\n", post.ID, post.UserID, post.Title)
	}
	if err := writer.Flush(); err != nil {
		return runtimeError(stderr, fmt.Errorf("write output: %w", err))
	}
	return ExitOK
}

func runGet(ctx context.Context, args []string, stdout, stderr io.Writer, client postsAPI) int {
	args, outputJSON := extractJSONFlag(args)
	flags := newFlagSet("get", stderr)
	if err := flags.Parse(args); err != nil {
		return flagErrorCode(err)
	}
	if flags.NArg() != 1 {
		return usageError(stderr, "get requires exactly one post ID")
	}

	id, err := positiveInt(flags.Arg(0), "post ID")
	if err != nil {
		return usageError(stderr, err.Error())
	}
	post, err := client.GetPost(ctx, id)
	if err != nil {
		return runtimeError(stderr, err)
	}
	if outputJSON {
		return writeJSON(stdout, stderr, post)
	}

	fmt.Fprintf(stdout, "ID: %d\nUser: %d\nTitle: %s\n\n%s\n", post.ID, post.UserID, post.Title, post.Body)
	return ExitOK
}

func runCreate(ctx context.Context, args []string, stdout, stderr io.Writer, client postsAPI) int {
	args, outputJSON := extractJSONFlag(args)
	flags := newFlagSet("create", stderr)
	title := flags.String("title", "", "post title (required)")
	body := flags.String("body", "", "post body (required)")
	userID := flags.Int("user-id", 0, "user ID (required)")
	if err := flags.Parse(args); err != nil {
		return flagErrorCode(err)
	}
	if flags.NArg() != 0 {
		return usageError(stderr, "create does not accept positional arguments")
	}
	if *title == "" {
		return usageError(stderr, "--title is required")
	}
	if *body == "" {
		return usageError(stderr, "--body is required")
	}
	if *userID < 1 {
		return usageError(stderr, "--user-id must be a positive integer")
	}

	// CLIの文字列・数値を、API層が受け取る構造体へ詰め替えます。
	post, err := client.CreatePost(ctx, api.CreatePostInput{
		UserID: *userID,
		Title:  *title,
		Body:   *body,
	})
	if err != nil {
		return runtimeError(stderr, err)
	}
	if outputJSON {
		return writeJSON(stdout, stderr, post)
	}

	fmt.Fprintf(stdout, "Created post %d\nTitle: %s\n", post.ID, post.Title)
	return ExitOK
}

func newFlagSet(name string, stderr io.Writer) *flag.FlagSet {
	// ContinueOnErrorを指定すると、flagがos.Exitせずエラーを返します。
	// これによりRunが終了コードを決められ、テストも途中終了しません。
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {}
	return flags
}

// extractJSONFlagは--jsonを引数列から取り出します。
// 標準のflagは最初の位置引数で解析を止めるため、この処理により
// 「get --json 1」と「get 1 --json」の両方を許可しています。
func extractJSONFlag(args []string) ([]string, bool) {
	filtered := make([]string, 0, len(args))
	outputJSON := false
	for _, arg := range args {
		if arg == "--json" {
			outputJSON = true
			continue
		}
		filtered = append(filtered, arg)
	}
	return filtered, outputJSON
}

func positiveInt(value, name string) (int, error) {
	// CLIの引数はすべて文字列なので、Atoiでintへ変換してから範囲を検証します。
	number, err := strconv.Atoi(value)
	if err != nil || number < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return number, nil
}

func writeJSON(stdout, stderr io.Writer, value any) int {
	// EncoderはJSONの末尾に改行も付けるため、CLIの出力に向いています。
	encoder := json.NewEncoder(stdout)
	// タイトルなどに含まれる<や>を、そのまま読みやすく出力します。
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return runtimeError(stderr, fmt.Errorf("write JSON: %w", err))
	}
	return ExitOK
}

func usageError(stderr io.Writer, message string) int {
	// 正常なデータはstdout、エラーはstderrへ分けるのがCLIの基本的な約束です。
	fmt.Fprintf(stderr, "error: %s\n", message)
	return ExitUsage
}

func runtimeError(stderr io.Writer, err error) int {
	fmt.Fprintf(stderr, "error: %v\n", err)
	return ExitRuntime
}

func flagErrorCode(err error) int {
	if errors.Is(err, flag.ErrHelp) {
		return ExitOK
	}
	return ExitUsage
}

func printUsage(writer io.Writer) {
	fmt.Fprintln(writer, `Usage:
  postcli list [--limit N] [--json]
  postcli get <id> [--json]
  postcli create --title TEXT --body TEXT --user-id ID [--json]
  postcli help`)
}
