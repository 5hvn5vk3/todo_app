package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

// /login - ログインフォーム表示
func login(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
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
			http.Redirect(w, r, "/todos", http.StatusFound)
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

// authenticate - ログイン認証処理
func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Password == models.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// logout - ログアウト処理
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")

	if err != nil {
		log.Println(err)
	}
	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID()
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
