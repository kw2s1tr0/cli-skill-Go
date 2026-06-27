package myapp

import (
	"aiagentcliapp/controller"
	"aiagentcliapp/repository"
	"context"
	"net/http"
	"os"
	"time"

	"golang.org/x/term"
)

// 環境変数をここで宣言する
const baseURL = "http://host.docker.internal"

func main() {
	// httpClientの共通設定をここで宣言する
	// タイムアウト設定
	httpClient := &http.Client{Timeout: 10 * time.Second}
	// 使用するclientインスタンスを作成する（参照のため使いまわす）
	client := repository.NewClient(baseURL, httpClient)

	// Controllerで引数ごとにserviceを呼び出し、結果をintで返す
	os.Exit(
		controller.Run(
			// 引数の最初は起動commandのため除外する
			os.Args[1:],
			client,
			context.Background(),
			os.Stdin,
			os.Stdout,
			os.Stderr,
			term.ReadPassword,
		),
	)
}
