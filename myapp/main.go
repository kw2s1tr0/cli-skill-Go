package myapp

import (
	"aiagentcliapp/controller"
	"aiagentcliapp/repository"
	"context"
	"net/http"
	"os"
	"time"
)

// 環境変数をここで宣言する
const baseURL = ""

func main() {
	// httpClientの共通設定をここで宣言する
	// タイムアウト設定
	httpClient := &http.Client{Timeout: 10 * time.Second}
	// 使用するclientインスタンスを作成する（参照のため使いまわす）
	client := repository.NewClient(baseURL, httpClient)

	// Controllerで引数ごとにserviceを呼び出し、結果をintで返す
	os.Exit(
		controller.Run(
			// 最初は起動commandのため除外する
			os.Args[1:],
			client,
			context.Background(),
			os.Stdout,
			os.Stderr,
		),
	)
}
