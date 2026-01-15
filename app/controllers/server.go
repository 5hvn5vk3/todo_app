package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"todo_app/app/models"
	"todo_app/config"
)

// generateHTML - HTMLテンプレートを生成
func generateHTML(w http.ResponseWriter, data any, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

// session - セッション情報を取得
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

// getMethod - HTTPメソッドを取得（_methodフィールドでのオーバーライドに対応）
func getMethod(r *http.Request) string {
	if r.Method == "POST" {
		if method := r.FormValue("_method"); method != "" {
			return method
		}
	}
	return r.Method
}

// StartMainServer - メインサーバーを起動
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
