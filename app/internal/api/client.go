package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// PostはJSONPlaceholderから返される投稿を表します。
// `json:"..."`は、JSONのキーとGoのフィールドを対応付ける構造体タグです。
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// CreatePostInputは投稿作成時にAPIへ送る項目だけを持ちます。
// 作成前にはIDが存在しないため、レスポンス用のPostとは型を分けています。
type CreatePostInput struct {
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Clientは投稿APIとのHTTP通信を担当します。
// フィールド名を小文字で始めると、このapiパッケージの外からは直接変更できません。
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClientはAPIクライアントを作ります。
// 接続先とHTTPクライアントを外から受け取ることで、テスト時はFake APIへ差し替えられます。
func NewClient(baseURL string, httpClient *http.Client) *Client {
	// nilでも動作するよう、Go標準のHTTPクライアントを予備として使います。
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		// 呼び出し側が末尾に/を付けてもURLが//にならないよう正規化します。
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

// ListPostsは最大limit件の投稿を取得します。
// context.Contextを渡すと、呼び出し側から通信のキャンセルや期限を伝えられます。
func (c *Client) ListPosts(ctx context.Context, limit int) ([]Post, error) {
	endpoint, err := url.Parse(c.baseURL + "/posts")
	if err != nil {
		return nil, fmt.Errorf("build list URL: %w", err)
	}

	query := endpoint.Query()
	// 数値をURLへ入れるには文字列へ変換します。Queryを使うと適切にエスケープされます。
	query.Set("_limit", strconv.Itoa(limit))
	endpoint.RawQuery = query.Encode()

	var posts []Post
	// &posts（ポインタ）を渡すことで、doJSONがデコード結果をpostsへ書き込めます。
	if err := c.doJSON(ctx, http.MethodGet, endpoint.String(), nil, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// GetPostはIDを指定して投稿を1件取得します。
func (c *Client) GetPost(ctx context.Context, id int) (Post, error) {
	var post Post
	endpoint := fmt.Sprintf("%s/posts/%d", c.baseURL, id)
	if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &post); err != nil {
		return Post{}, err
	}
	return post, nil
}

// CreatePostは新しい投稿をJSONで送り、APIから返された投稿を取得します。
func (c *Client) CreatePost(ctx context.Context, input CreatePostInput) (Post, error) {
	var post Post
	if err := c.doJSON(ctx, http.MethodPost, c.baseURL+"/posts", input, &post); err != nil {
		return Post{}, err
	}
	return post, nil
}

func (c *Client) doJSON(ctx context.Context, method, endpoint string, body any, result any) error {
	// GETには本文がありません。POSTなどbodyがある場合だけJSONへ変換します。
	// anyは「任意の型」を受け取れるGoの型です。
	var requestBody io.Reader
	if body != nil {
		var buffer bytes.Buffer
		if err := json.NewEncoder(&buffer).Encode(body); err != nil {
			return fmt.Errorf("encode request body: %w", err)
		}
		requestBody = &buffer
	}

	request, err := http.NewRequestWithContext(ctx, method, endpoint, requestBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("request API: %w", err)
	}
	// HTTPレスポンスの本文は必ず閉じます。deferはこの関数を抜ける直前に実行されます。
	defer response.Body.Close()

	// HTTPは通信自体に成功しても404や500を返すため、ステータスコードを別途確認します。
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		// エラー本文が巨大でもメモリを使いすぎないよう、先頭1KiBだけ読みます。
		message, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		if text := strings.TrimSpace(string(message)); text != "" {
			return fmt.Errorf("API returned %s: %s", response.Status, text)
		}
		return fmt.Errorf("API returned %s", response.Status)
	}

	// resultは呼び出し元から渡されたポインタです。JSONをその参照先へ書き込みます。
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return fmt.Errorf("decode API response: %w", err)
	}
	return nil
}
