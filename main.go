package main

import (
	"log"
	"todo_app/app/controllers"
	"todo_app/app/models"
)

func main() {
	// データベース接続の確認
	if models.Db == nil {
		log.Fatal("Database connection is not initialized")
	}
	log.Println("Database instance:", models.Db)

	// サーバー起動
	log.Println("Starting server...")
	if err := controllers.StartMainServer(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
