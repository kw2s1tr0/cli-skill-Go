package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(
	baseURL string,
	httpClient *http.Client,
) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

// API呼び出し共通関数
func (client *Client) DoJSON(
	ctx context.Context,
	method string,
	path string,
	body any,
	result any,
) error {

	// io.Reader HTTPリクエストの本文の型
	var requestBody io.Reader

	// bodyがあるとき（get以外）はJSONに変換
	if body != nil {

		// bytes.Bufferは一時的に入れる箱の型
		var buffer bytes.Buffer
		if err := json.NewEncoder(&buffer).Encode(body); err != nil {
			return fmt.Errorf("encode request body: %w", err)
		}
		requestBody = &buffer
	}

	// HTTPリクエストを作成
	request, err := http.NewRequestWithContext(ctx, method, client.baseURL+path, requestBody)

	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	request.Header.Set("Accept", "application/json")

	// get以外はJSONで送る
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	// HTTP通信
	response, err := client.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("request API: %w", err)
	}

	// deferは関数の最後に実行されるようにする記法
	// bodyを閉じて、通信自体を閉じる
	defer response.Body.Close()

	// 異常値返却
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {

		// 本文を1024バイト読む 巨大レスポンスでのメモリ圧迫を避ける
		message, _ := io.ReadAll(io.LimitReader(response.Body, 1024))

		// 本文が空でなければそれを返し、空ならステータスのみ
		if text := formatAPIErrorMessage(message); text != "" {
			return fmt.Errorf("API returned %s: %s", response.Status, text)
		}
		return fmt.Errorf("API returned %s", response.Status)
	}

	if result == nil {
		return nil
	}

	// Jsonからオブジェクトに変換
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return fmt.Errorf("decode API response: %w", err)
	}

	// 参照されたresultに値を詰めているため結果としては返却しない
	return nil
}

func formatAPIErrorMessage(message []byte) string {
	text := strings.TrimSpace(string(message))
	if text == "" {
		return ""
	}

	var apiError struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(message, &apiError); err != nil {
		return text
	}

	if decodedMessage := strings.TrimSpace(apiError.Message); decodedMessage != "" {
		return decodedMessage
	}
	return text
}
