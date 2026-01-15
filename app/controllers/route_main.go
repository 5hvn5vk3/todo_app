package controllers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"todo_app/app/models"
)

// top - トップページ表示
func top(w http.ResponseWriter, r *http.Request) {
	// "/" のみを処理
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, err := session(w, r)
	if err != nil {
		generateHTML(w, "Hello", "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", http.StatusSeeOther)
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

// index - Todo一覧表示
func index(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	todos, err := user.GetTodosByUser()
	if err != nil {
		log.Println(err)
	}
	user.Todos = todos
	generateHTML(w, user, "layout", "private_navbar", "index")
}

// todoNew - Todo新規作成フォーム表示
func todoNew(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

// todoSave - Todo保存処理
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

// todoEdit - Todo編集フォーム表示
func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err = sess.GetUserBySession()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	t, err := models.GetTodo(id)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}
	generateHTML(w, t, "layout", "private_navbar", "todo_edit")
}

// todoUpdate - Todo更新処理
func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
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
	t := &models.Todo{
		ID:      id,
		Content: content,
		UserID:  user.ID,
	}
	if err := t.UpdateTodo(); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}

// todoDelete - Todo削除処理
func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err = sess.GetUserBySession()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	t, err := models.GetTodo(id)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	if err := t.DeleteTodo(); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}
