# Minerva Web

書籍情報管理システムのWeb UIです。

## 技術スタック

- **Backend**: Go 1.26 (net/http)
- **Frontend**: HTMX 2.0 + water.css
- **Database**: PostgreSQL 16

## 機能

- 書籍情報のCRUD操作
- リアルタイム検索（タイトル・著者・ISBN）
- 列の並び替え
- ページネーション（1ページあたりの件数変更可能）
- ソフトデリート（deleted_atによる論理削除）

## 起動方法

### Docker Compose（推奨）

```bash
# .env ファイルを確認/編集
cat .env

# 起動
docker compose up -d

# ログ確認
docker compose logs -f app
```

アプリケーションは http://localhost:8081 でアクセス可能です。

### ローカル開発

```bash
# 依存パッケージのインストール
go mod download

# ビルド
go build -o minerva-web .

# 実行
./minerva-web
```

## 環境変数

`.env` ファイルで以下を設定します：

```
POSTGRES_DB=bibliography
POSTGRES_USER=bibliography
POSTGRES_PASSWORD=bibliography
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

## ディレクトリ構成

```
.
├── main.go                    # エントリポイント
├── internal/
│   ├── config/config.go       # 設定読み込み
│   ├── db/db.go               # DB接続プール
│   ├── model/book.go          # Bookモデル
│   ├── repository/book.go     # DB操作
│   └── handler/
│       ├── handler.go         # ルーティング、テンプレート
│       ├── page.go            # ページ表示
│       ├── book.go            # CRUDハンドラ
│       └── handler_test.go    # テスト
├── templates/                 # HTMLテンプレート
├── static/                    # 静的ファイル（CSS, JS）
└── schema.sql                 # DBスキーマ
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

## ライセンス

Private
