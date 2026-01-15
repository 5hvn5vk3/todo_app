# Todo アプリケーション

## 概要

[Udemy の講座](https://www.udemy.com/share/103TVa3@JcVxm1FpsGbhhZMFmgga_7FYkIvuEsPq0odEDP73o43qRc4NIgcXO02wTqP5wS5IRA==/)を見ながら作成した Todo アプリケーションを、より RESTful な設計に準拠するようにリファクタリングしました。

### 本アプリケーションの機能

- ユーザー登録・ログイン・ログアウト
- Todo の作成・閲覧・更新・削除（CRUD）
- セッション管理によるアクセス制御

### 主な改善点

1. **動詞ベースの URL からリソース指向の URL へ移行**
2. **適切な HTTP ステータスコードの使用**
3. **HTTP メソッドオーバーライド機能の実装**
4. **エラーハンドリングの強化**
5. **ルーティングの整理と最適化**

[リファクタリングの詳細はこちら](https://github.com/5hvn5vk3/todo_app/blob/render-rest/RESTful%E8%A8%AD%E8%A8%88%E3%81%B8%E3%81%AE%E6%94%B9%E5%96%84.md)

---

## デプロイ先

[Render にデプロイしています。](https://todo-app-czkz.onrender.com)  
※起動に少々時間がかかる場合があります。

---

## アプリケーション設計

### 技術スタック

- **言語**: Go 1.25.0
- **データベース**: PostgreSQL（本番環境）/ SQLite3（開発環境）
- **フロントエンド**: Bootstrap 5, jQuery 3.7.1
- **テンプレートエンジン**: html/template
- **認証**: セッションベース（Cookie）
- **デプロイ**: Render

### アーキテクチャ

Go 言語による MVC パターンに基づいた Web アプリケーションで、以下の 3 層構造で構成されています：

- **Models（モデル層）**: データベースとのやり取りとビジネスロジック
- **Views（ビュー層）**: HTML テンプレートによる UI 表示
- **Controllers（コントローラ層）**: ルーティングと HTTP リクエスト処理

### ディレクトリ構造

```
todo_app/
├── main.go                 # エントリーポイント
├── go.mod                  # Go モジュール定義
├── config.ini              # 設定ファイル
├── webapp.sql              # データベーススキーマ
├── app/
│   ├── controllers/        # コントローラ層
│   │   ├── server.go       # サーバー設定とミドルウェア
│   │   ├── route_main.go   # メインルーティング（Todo CRUD）
│   │   └── route_auth.go   # 認証ルーティング（ログイン/サインアップ）
│   ├── models/             # モデル層
│   │   ├── base.go         # データベース接続
│   │   ├── users.go        # ユーザーモデル
│   │   └── todo.go         # Todoモデル
│   └── views/              # ビュー層
│       ├── templates/      # HTMLテンプレート
│       ├── css/            # スタイルシート
│       └── js/             # JavaScriptファイル
├── config/
│   └── config.go           # 設定読み込み
└── utils/
    └── logging.go          # ロギング機能
```

---
