# RESTful 設計への改善

## 概要

[Udemy の講座](https://www.udemy.com/share/103TVa3@JcVxm1FpsGbhhZMFmgga_7FYkIvuEsPq0odEDP73o43qRc4NIgcXO02wTqP5wS5IRA==/)を見ながら作成した Todo アプリケーションを、より RESTful な設計に準拠するようにリファクタリングしました。主な改善点は以下の通りです：

1. [**RESTful ルーティングへの移行（リソース指向 URL とハンドラー実装）**](#1-restful-ルーティングへの移行)
2. [**HTTP メソッドオーバーライド機能の実装**](#2-http-メソッドオーバーライド機能の実装)
3. [**適切な HTTP ステータスコードの使用**](#3-http-ステータスコードの改善)

---

## 1. RESTful ルーティングへの移行

### 対象ファイル

- `app/controllers/server.go` - ルーティング定義
- `app/controllers/route_auth.go` - 認証関連のハンドラー
- `app/controllers/route_main.go` - Todo 関連のハンドラー
- `app/views/templates/*.html` - フォームの action 属性

### URL 設計の改善（動詞ベースからリソース指向へ）

#### 変更前（動詞ベース URL）

```
GET  /signup          → サインアップフォーム
POST /signup          → ユーザー作成
POST /authenticate    → ログイン
GET  /logout          → ログアウト
GET  /todos           → Todo一覧
GET  /todos/new       → Todo作成フォーム
POST /todos/save      → Todo作成
GET  /todos/edit/:id  → Todo編集フォーム
POST /todos/update/:id → Todo更新
GET  /todos/delete/:id → Todo削除
```

#### 変更後（リソース指向 URL）

```
GET    /users           → サインアップフォーム
POST   /users           → ユーザー作成
POST   /sessions        → ログイン（セッション作成）
DELETE /sessions        → ログアウト（セッション削除）
GET    /todos           → Todo一覧
POST   /todos           → Todo作成
GET    /todos/new       → Todo作成フォーム
GET    /todos/:id       → Todo編集フォーム
PUT    /todos/:id       → Todo更新
DELETE /todos/:id       → Todo削除
```

#### 改善のポイント

- **リソース名は名詞（複数形）**：`/users`, `/sessions`, `/todos`
- **HTTP メソッドで操作を表現**：GET（取得）、POST（作成）、PUT（更新）、DELETE（削除）
- **URL から動詞を排除**：`/authenticate`, `/logout`, `/save`, `/update`, `/delete` などを廃止
- **階層的なリソース表現**：`/todos/:id` で Todo リソースの個別操作

### ルーティング構造の実装

#### `app/controllers/server.go`

```go
func StartMainServer() error {
    files := http.FileServer(http.Dir(config.Config.Static))
    http.Handle("/static/", http.StripPrefix("/static/", files))

    // Public routes
    http.HandleFunc("/", top)
    http.HandleFunc("/login", login)

    // RESTful routes
    http.HandleFunc("/users", users)       // GET: サインアップフォーム, POST: ユーザー作成
    http.HandleFunc("/sessions", sessions) // POST: ログイン, DELETE: ログアウト

    // Todos routes (RESTful)
    http.HandleFunc("/todos/", todosHandler)     // GET: 編集フォーム, PUT: 更新, DELETE: 削除 (/todos/:id)
    http.HandleFunc("/todos", todosIndexHandler) // GET: 一覧, POST: 作成

    port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, nil)
}
```

### リソースごとのハンドラー実装

#### `/users` - ユーザーリソース（`app/controllers/route_auth.go`）

```go
func users(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        // サインアップフォーム表示
    case "POST":
        // ユーザー作成
    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
}
```

#### `/sessions` - セッションリソース（`app/controllers/route_auth.go`）

```go
func sessions(w http.ResponseWriter, r *http.Request) {
    method := getMethod(r)
    switch method {
    case "POST":
        // ログイン（セッション作成）
        authenticate(w, r)
    case "DELETE":
        // ログアウト（セッション削除）
        logout(w, r)
    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
}
```

#### `/todos` - Todo リストリソース（`app/controllers/route_main.go`）

```go
func todosIndexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/todos" {
        return // 次のハンドラーに処理を委譲
    }

    switch r.Method {
    case "GET":
        // Todo一覧表示
        index(w, r)
    case "POST":
        // Todo作成
        todoSave(w, r)
    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
}
```

#### `/todos/:id` - 個別 Todo リソース（`app/controllers/route_main.go`）

```go
func todosHandler(w http.ResponseWriter, r *http.Request) {
    validPath := regexp.MustCompile("^/todos/([0-9]+)$")
    matches := validPath.FindStringSubmatch(r.URL.Path)

    if matches == nil {
        if r.URL.Path == "/todos/new" {
            todoNew(w, r)
            return
        }
        http.NotFound(w, r)
        return
    }

    id, err := strconv.Atoi(matches[1])
    if err != nil {
        http.NotFound(w, r)
        return
    }

    method := getMethod(r)
    switch method {
    case "GET":
        // Todo編集フォーム表示
        todoEdit(w, r, id)
    case "PUT", "PATCH":
        // Todo更新
        todoUpdate(w, r, id)
    case "DELETE":
        // Todo削除
        todoDelete(w, r, id)
    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
}
```

---

## 2. HTTP メソッドオーバーライド機能の実装

HTML フォームは `GET` と `POST` しかサポートしていないため、`PUT` と `DELETE` を実現するための機能を実装しました。

### 実装内容

#### 対象ファイル: `app/controllers/server.go`

```go
// HTTPメソッドを取得（_methodフィールドでのオーバーライドに対応）
func getMethod(r *http.Request) string {
    if r.Method == "POST" {
        if method := r.FormValue("_method"); method != "" {
            return method
        }
    }
    return r.Method
}
```

### 使用例（HTML テンプレート）

#### DELETE（Todo 削除） - `app/views/templates/index.html`

```html
<form action="/todos/{{.ID}}" method="post" style="display:inline;">
  <input type="hidden" name="_method" value="DELETE" />
  <button type="submit">[Delete]</button>
</form>
```

#### PUT（Todo 更新） - `app/views/templates/todo_edit.html`

```html
<form role="form" action="/todos/{{.ID}}" method="post">
  <input type="hidden" name="_method" value="PUT" />
  <textarea class="form-control" name="content">{{.Content}}</textarea>
  <button type="submit">Update</button>
</form>
```

#### DELETE（ログアウト） - `app/views/templates/private_navbar.html`

```html
<form action="/sessions" method="post" style="display:inline;">
  <input type="hidden" name="_method" value="DELETE" />
  <button type="submit">logout</button>
</form>
```

---

## 3. HTTP ステータスコードの改善

### リダイレクトのステータスコード

#### 対象ファイル: `app/controllers/route_auth.go`, `app/controllers/route_main.go`

**変更前**：全てのリダイレクトで `302 Found` を使用

```go
http.Redirect(w, r, "/todos", http.StatusFound) // 302
```

**変更後**：POST 後のリダイレクトには `303 See Other` を使用

```go
http.Redirect(w, r, "/todos", http.StatusSeeOther) // 303
```

**理由**：

- `303 See Other` は、POST 後のリダイレクトで次のリクエストを必ず GET にする
- ブラウザの「戻る」ボタンや「再読み込み」で二重送信（Double Submit）を防止
- RESTful 設計のベストプラクティスに準拠

### エラーハンドリングのステータスコード

#### 対象ファイル: `app/controllers/route_auth.go`, `app/controllers/route_main.go`

**変更前**：エラー時もログ出力のみで適切な HTTP ステータスを返していなかった

```go
err := r.ParseForm()
if err != nil {
    log.Println(err)
    // ステータスコードなし
}
```

**変更後**：適切な HTTP ステータスコードを返すように改善

```go
err := r.ParseForm()
if err != nil {
    log.Println(err)
    http.Error(w, "Bad Request", http.StatusBadRequest) // 400
    return
}
```

#### 使用するステータスコード一覧

| コード                      | 説明                           | 使用場面                                   |
| --------------------------- | ------------------------------ | ------------------------------------------ |
| `303 See Other`             | POST 後のリダイレクト          | ユーザー作成、Todo 作成・更新・削除後      |
| `400 Bad Request`           | 不正なリクエスト               | フォームのパースエラー                     |
| `401 Unauthorized`          | 認証エラー                     | セッションが無効、ユーザー取得失敗         |
| `404 Not Found`             | リソースが存在しない           | Todo が見つからない                        |
| `405 Method Not Allowed`    | 許可されていない HTTP メソッド | サポートされていないメソッドでのリクエスト |
| `500 Internal Server Error` | サーバー内部エラー             | データベースエラー、セッション作成失敗     |

---

## 4. その他の改善

### エラーハンドリングの強化

#### 改善前の問題点

- エラーが発生してもログ出力のみで処理を継続
- クライアントに適切なエラーレスポンスを返していない
- `else` ブロックの深いネスト

#### 改善後のエラーハンドリング

**例: todoSave 関数（`app/controllers/route_main.go`）**

```go
func todoSave(w http.ResponseWriter, r *http.Request) {
    sess, err := session(w, r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    err = r.ParseForm()
    if err != nil {
        log.Println(err)
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    user, err := sess.GetUserBySession()
    if err != nil {
        log.Println(err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    content := r.PostFormValue("content")
    if err := user.CreateTodo(content); err != nil {
        log.Println(err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/todos", http.StatusSeeOther)
}
```

**改善ポイント：**

1. **Early Return パターン**：エラーが発生したら即座に return
2. **適切なステータスコード**：エラーの種類に応じた HTTP ステータスを返す
3. **ネストの削減**：`else` ブロックを排除し、コードの可読性を向上

### トップページのルーティング改善（`app/controllers/route_main.go`）

```go
func top(w http.ResponseWriter, r *http.Request) {
    // "/" のみを処理（他のパスは404）
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    // ... 処理 ...
}
```

### ログ出力の追加（`main.go`）

```go
func main() {
    fmt.Println(models.Db)

    log.Println("Starting server on port:", "8080")
    if err := controllers.StartMainServer(); err != nil {
        log.Fatal(err)
    }
}
```

---

## 変更されたファイル一覧

| ファイル                                  | 変更内容                                                                                           |
| ----------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `app/controllers/server.go`               | ルーティングの再設計、HTTP メソッドオーバーライド機能の実装                                        |
| `app/controllers/route_auth.go`           | 認証関連ハンドラーの実装、ステータスコード改善、エラーハンドリング強化                             |
| `app/controllers/route_main.go`           | Todo 関連ハンドラーの実装、ステータスコード改善、エラーハンドリング強化、Early Return パターン適用 |
| `app/views/templates/index.html`          | RESTful URL に対応、DELETE メソッド実装                                                            |
| `app/views/templates/login.html`          | `/sessions`エンドポイントに変更                                                                    |
| `app/views/templates/signup.html`         | `/users`エンドポイントに変更                                                                       |
| `app/views/templates/todo_edit.html`      | RESTful URL に対応、PUT メソッド実装                                                               |
| `app/views/templates/todo_new.html`       | RESTful URL に対応                                                                                 |
| `app/views/templates/private_navbar.html` | DELETE メソッドでログアウト実装                                                                    |
| `app/views/templates/public_navbar.html`  | サインアップリンクを`/users`エンドポイントに変更                                                   |
| `main.go`                                 | ログ出力とエラーハンドリング追加                                                                   |
