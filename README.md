# cli-skill-Go

Go製CLIを作り、将来Codex Skillから自作アプリを操作するための練習プロジェクトです。
最初の題材として、公開Fake APIの
[JSONPlaceholder](https://jsonplaceholder.typicode.com/)から投稿を取得・作成する
`postcli`を実装しています。

## 実行方法

Go moduleは`app`ディレクトリにあります。

```bash
cd app
go run ./cmd/postcli help
go run ./cmd/postcli list --limit 3
go run ./cmd/postcli get 1
go run ./cmd/postcli get 1 --json
go run ./cmd/postcli create \
  --title "Go CLIの練習" \
  --body "JSONPlaceholderへ送信します" \
  --user-id 1 \
  --json
```

単体の実行ファイルも作れます。

```bash
cd app
go build -o postcli ./cmd/postcli
./postcli list --limit 3
```

終了コードは、成功が`0`、通信・APIエラーが`1`、引数エラーが`2`です。
通常は人向けの表示を出し、`--json`を付けるとAIから扱いやすいJSONを標準出力へ返します。
エラーは標準エラー出力へ返します。

JSONPlaceholderのPOSTは練習用の擬似処理です。作成レスポンスは返りますが、データは永続化されません。

## コード構成

```text
app/
├── cmd/postcli/main.go       # CLIの起動、HTTPタイムアウト、終了
├── internal/cli/run.go       # 引数解析、表示、終了コード
└── internal/api/client.go    # HTTPリクエスト、JSON変換
```

`main`からCLI処理とAPI通信を分離することで、公開APIへ接続せずに各層をテストできます。

## 検証

```bash
cd app
go test ./...
go vet ./...
go build ./cmd/postcli
```
