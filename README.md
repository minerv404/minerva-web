# minerva-web

書籍管理用 Web フロントエンド + バックエンド API

## 構成

- バックエンド: Go 1.26 (net/http)
- フロントエンド: HTMX 2.0 + water.css
- データベース: PostgreSQL 16

## 起動方法

### Docker Compose

```bash
vim .env
docker compose up -d
```

http://localhost:8081 で起動します。

### 手元での開発

```bash
go mod download
go build -o minerva-web .
./minerva-web
```

## 環境変数

`.env` ファイルで以下を設定します。

```
POSTGRES_DB=<your-db-host>
POSTGRES_USER=<your-db-user>
POSTGRES_PASSWORD=<your-db-password>
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```
## テスト

```bash
go test ./...
```

## API エンドポイント

| Method | Path | 説明 |
|--------|------|------|
| GET | `/` | 書籍一覧ページ |
| GET | `/books?q=&sort=&order=&per_page=&page=` | 書籍一覧（検索・ソート・ページネーション） |
| GET | `/books/new` | 新規登録フォーム |
| GET | `/books/{isbn}` | 書籍詳細 |
| POST | `/books` | 書籍作成 |
| GET | `/books/{isbn}/edit` | 編集フォーム |
| PUT | `/books/{isbn}` | 書籍更新 |
| DELETE | `/books/{isbn}` | 書籍削除（ソフトデリート） |
