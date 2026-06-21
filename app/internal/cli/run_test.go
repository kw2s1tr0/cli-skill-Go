package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"company-cli/internal/api"
)

// fakeAPIは、実際のHTTP通信をせずにCLIをテストするための代役（Test Double）です。
// 各テストが用意した戻り値をそのまま返すため、CLIの引数解析と表示だけに集中できます。
type fakeAPI struct {
	posts       []api.Post
	post        api.Post
	createdPost api.Post
	err         error
}

// この代入は実行目的ではなく、fakeAPIがpostsAPIを満たすことをコンパイル時に確認します。
// postsAPIへメソッドを追加したのにFakeを直し忘れると、ここでエラーになります。
var _ postsAPI = fakeAPI{}

func (fake fakeAPI) ListPosts(context.Context, int) ([]api.Post, error) {
	return fake.posts, fake.err
}

func (fake fakeAPI) GetPost(context.Context, int) (api.Post, error) {
	return fake.post, fake.err
}

func (fake fakeAPI) CreatePost(context.Context, api.CreatePostInput) (api.Post, error) {
	return fake.createdPost, fake.err
}

func TestListHumanOutput(t *testing.T) {
	// Arrange: APIが返す投稿を準備します。
	client := fakeAPI{posts: []api.Post{{ID: 1, UserID: 2, Title: "hello"}}}
	// Act: 実際のコマンドラインと同じ形の引数でRunを呼びます。
	code, stdout, stderr := runForTest([]string{"list", "--limit", "1"}, client)

	// Assert: 終了コードと、人向け表示の重要な部分を確認します。
	// 空白の数まで固定すると表示の小変更で壊れやすいため、必要な文字列だけを確認します。
	if code != ExitOK {
		t.Fatalf("exit code = %d, want %d; stderr = %s", code, ExitOK, stderr)
	}
	if !strings.Contains(stdout, "ID") || !strings.Contains(stdout, "hello") {
		t.Fatalf("stdout = %q, want table header and post title", stdout)
	}
}

func TestGetJSONAfterID(t *testing.T) {
	// --jsonが位置引数（7）の後ろにあっても認識され、JSONを返すことを確認します。
	client := fakeAPI{post: api.Post{ID: 7, UserID: 3, Title: "seven", Body: "body"}}
	code, stdout, stderr := runForTest([]string{"get", "7", "--json"}, client)

	if code != ExitOK {
		t.Fatalf("exit code = %d, want %d; stderr = %s", code, ExitOK, stderr)
	}
	if want := `"id":7`; !strings.Contains(stdout, want) {
		t.Fatalf("stdout = %q, want it to contain %s", stdout, want)
	}
}

func TestCreateRequiresTitle(t *testing.T) {
	// 必須フラグがない場合はAPIを呼ばず、使い方のエラー（終了コード2）にします。
	code, _, stderr := runForTest([]string{"create", "--body", "body", "--user-id", "1"}, fakeAPI{})

	if code != ExitUsage {
		t.Fatalf("exit code = %d, want %d", code, ExitUsage)
	}
	if !strings.Contains(stderr, "--title is required") {
		t.Fatalf("stderr = %q, want missing title error", stderr)
	}
}

func TestInvalidLimit(t *testing.T) {
	// flagパッケージがintへ変換できても、業務上許可しない範囲はCLI側で検証します。
	code, _, stderr := runForTest([]string{"list", "--limit", "0"}, fakeAPI{})

	if code != ExitUsage || !strings.Contains(stderr, "between 1 and 100") {
		t.Fatalf("code = %d, stderr = %q, want usage error for limit", code, stderr)
	}
}

func TestUnknownCommand(t *testing.T) {
	// 未知のコマンドが成功扱いにならず、利用者へ理由を伝えることを確認します。
	code, _, stderr := runForTest([]string{"unknown"}, fakeAPI{})

	if code != ExitUsage || !strings.Contains(stderr, "unknown command") {
		t.Fatalf("code = %d, stderr = %q, want unknown command usage error", code, stderr)
	}
}

func TestRuntimeError(t *testing.T) {
	// API層の失敗をFakeで発生させ、stderrと実行エラー（終了コード1）を確認します。
	code, _, stderr := runForTest([]string{"get", "1"}, fakeAPI{err: errors.New("network down")})

	if code != ExitRuntime || !strings.Contains(stderr, "network down") {
		t.Fatalf("code = %d, stderr = %q, want runtime error", code, stderr)
	}
}

func runForTest(args []string, client postsAPI) (int, string, string) {
	// bytes.Bufferはメモリ上のio.Writerです。
	// os.Stdout/os.Stderrの代わりに渡すと、出力内容を文字列として検査できます。
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run(context.Background(), args, &stdout, &stderr, client)
	return code, stdout.String(), stderr.String()
}
