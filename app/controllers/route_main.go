package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

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

func todoNew(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

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
