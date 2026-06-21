package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"company-cli/internal/api"
	"company-cli/internal/cli"
)

const baseURL = "https://jsonplaceholder.typicode.com"

func main() {
	// mainの仕事は「必要な部品を作ってCLIを起動する」だけにします。
	// 実際の処理を別パッケージへ分けると、mainを起動せずにテストできます。

	// HTTP通信が永久に待ち続けないよう、クライアント全体にタイムアウトを設定します。
	httpClient := &http.Client{Timeout: 10 * time.Second}
	client := api.NewClient(baseURL, httpClient)

	// os.Args[0]は実行ファイル名なので、CLIへ渡すのは[1:]以降です。
	// stdoutとstderrを分けると、AIやシェルは正常な出力だけを安全に利用できます。
	// os.ExitへRunの戻り値を渡し、成功・実行エラー・入力エラーを終了コードで伝えます。
	os.Exit(cli.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr, client))
}
