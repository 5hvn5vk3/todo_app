package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

func siginup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		generateHTML(w, nil, "layout", "public_navbar", "signup")
	case "POST":
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		if err := user.CreateUser(); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusFound) // 302
	}
}
