package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Goは「_test.go」で終わるファイルを、go test実行時だけコンパイルします。
// Testから始まり、引数に*testing.Tを取る関数が1つのテストとして実行されます。
func TestListPosts(t *testing.T) {
	// このテストは他のt.Parallel付きテストと並行実行できます。
	// テスト同士で共有するデータがある場合は、競合を避けるため使用できません。
	t.Parallel()

	// httptest.NewServerは、テスト中だけ使えるローカルHTTPサーバーを起動します。
	// 公開APIへ接続しないため、ネットワークや外部サービスの状態に左右されません。
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Clientが送ったリクエストをサーバー側で検査します。
		if request.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", request.Method)
		}
		if request.URL.Path != "/posts" {
			t.Errorf("path = %s, want /posts", request.URL.Path)
		}
		if got := request.URL.Query().Get("_limit"); got != "2" {
			t.Errorf("_limit = %q, want 2", got)
		}
		writer.Header().Set("Content-Type", "application/json")
		// 検査後、Fake APIとして決められたJSONを返します。
		_, _ = writer.Write([]byte(`[{"userId":1,"id":1,"title":"one","body":"body"}]`))
	}))
	// deferにより、成功・失敗のどちらでもテスト終了時にサーバーを停止します。
	defer server.Close()

	// Arrange（準備）は上まで、ここがAct（テスト対象の実行）です。
	posts, err := NewClient(server.URL, server.Client()).ListPosts(context.Background(), 2)
	// ここからAssert（結果の確認）です。
	// Fatalfはこのテストを即座に終了するため、以降を確認できないエラーに使います。
	if err != nil {
		t.Fatalf("ListPosts() error = %v", err)
	}
	if len(posts) != 1 || posts[0].Title != "one" {
		t.Fatalf("ListPosts() = %#v, want one post titled one", posts)
	}
}

func TestGetPost(t *testing.T) {
	t.Parallel()

	// IDがURLのパスへ正しく埋め込まれることと、JSONからPostへの変換を確認します。
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/posts/7" {
			t.Errorf("path = %s, want /posts/7", request.URL.Path)
		}
		_, _ = writer.Write([]byte(`{"userId":3,"id":7,"title":"seven","body":"body"}`))
	}))
	defer server.Close()

	post, err := NewClient(server.URL, server.Client()).GetPost(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}
	if post.ID != 7 || post.UserID != 3 {
		t.Fatalf("GetPost() = %#v, want ID 7 and user ID 3", post)
	}
}

func TestCreatePost(t *testing.T) {
	t.Parallel()

	// POSTではメソッドとヘッダーに加えて、送信されたJSON本文も検査します。
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", request.Method)
		}
		if contentType := request.Header.Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", contentType)
		}

		var input CreatePostInput
		// リクエスト本文を構造体へ戻し、期待した値が送信されたか確認します。
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if input.Title != "new post" || input.UserID != 4 {
			t.Errorf("input = %#v, want title new post and user ID 4", input)
		}
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"userId":4,"id":101,"title":"new post","body":"hello"}`))
	}))
	defer server.Close()

	post, err := NewClient(server.URL, server.Client()).CreatePost(context.Background(), CreatePostInput{
		UserID: 4,
		Title:  "new post",
		Body:   "hello",
	})
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}
	if post.ID != 101 {
		t.Fatalf("CreatePost().ID = %d, want 101", post.ID)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	// 正常系だけでなく、HTTP 404がGoのerrorとして返ることも確認します。
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		http.Error(writer, "missing post", http.StatusNotFound)
	}))
	defer server.Close()

	_, err := NewClient(server.URL, server.Client()).GetPost(context.Background(), 999)
	if err == nil || !strings.Contains(err.Error(), "404 Not Found") {
		t.Fatalf("GetPost() error = %v, want a 404 error", err)
	}
}

func TestInvalidJSON(t *testing.T) {
	t.Parallel()

	// HTTP 200でも本文が壊れている可能性があります。JSON解析失敗の経路を再現します。
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = writer.Write([]byte(`not-json`))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, server.Client()).GetPost(context.Background(), 1)
	if err == nil || !strings.Contains(err.Error(), "decode API response") {
		t.Fatalf("GetPost() error = %v, want a decode error", err)
	}
}

func TestTimeout(t *testing.T) {
	t.Parallel()

	// サーバーを意図的に遅らせ、HTTPクライアントのタイムアウトを確実に発生させます。
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_, _ = writer.Write([]byte(`[]`))
	}))
	defer server.Close()

	httpClient := server.Client()
	// サーバーの50msより短い1msを指定しているため、リクエストは失敗するはずです。
	httpClient.Timeout = time.Millisecond
	_, err := NewClient(server.URL, httpClient).ListPosts(context.Background(), 1)
	if err == nil || !strings.Contains(err.Error(), "request API") {
		t.Fatalf("ListPosts() error = %v, want a request timeout", err)
	}
}
