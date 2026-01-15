package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"todo_app/app/models"
	"todo_app/config"
)

func generateHTML(w http.ResponseWriter, data any, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid Session")
		}
	}
	return sess, err
}

// HTTPメソッドを取得（_methodフィールドでのオーバーライドに対応）
func getMethod(r *http.Request) string {
	if r.Method == "POST" {
		if method := r.FormValue("_method"); method != "" {
			return method
		}
	}
	return r.Method
}

// /users - ユーザー登録（RESTful）
func users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// サインアップフォーム表示
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			http.Redirect(w, r, "/todos", http.StatusSeeOther)
		}
	case "POST":
		// ユーザー作成
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		if err := user.CreateUser(); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// 201 Created + リダイレクト
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// /sessions - セッション管理（RESTful）
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

// /todos - Todoリスト取得・作成（RESTful）
func todosIndexHandler(w http.ResponseWriter, r *http.Request) {
	// パスが完全に /todos の場合のみ処理
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

// /todos/:id - 個別Todo操作（RESTful）
func todosHandler(w http.ResponseWriter, r *http.Request) {
	// /todos/ で始まり、その後に数字が続くパスをマッチ
	validPath := regexp.MustCompile("^/todos/([0-9]+)$")
	matches := validPath.FindStringSubmatch(r.URL.Path)

	if matches == nil {
		// /todos/new は新規作成フォーム
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

func StartMainServer() error {
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	// Public routes
	http.HandleFunc("/", top)
	http.HandleFunc("/login", login) // GET: ログインフォーム

	// RESTful routes
	http.HandleFunc("/users", users)       // GET: サインアップフォーム, POST: ユーザー作成
	http.HandleFunc("/sessions", sessions) // POST: ログイン, DELETE: ログアウト

	// Todos routes (RESTful)
	http.HandleFunc("/todos/", todosHandler)     // GET: 編集フォーム, PUT: 更新, DELETE: 削除 (/todos/:id)
	http.HandleFunc("/todos", todosIndexHandler) // GET: 一覧, POST: 作成

	port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, nil)
}
