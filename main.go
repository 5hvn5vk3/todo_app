package main

import (
	"fmt"
	"log"
	"todo_app/app/controllers"
	"todo_app/app/models"
)

func main() {
	fmt.Println(models.Db)

	log.Println("Starting server on port:", "8080")
	if err := controllers.StartMainServer(); err != nil {
		log.Fatal(err)
	}
}
