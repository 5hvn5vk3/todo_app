package main

import (
	"log"
	"todo_app/app/models"
	_ "todo_app/config"
)

func main() {
	// log.Println(config.Config.Port)
	// log.Println(config.Config.SQLDriver)
	// log.Println(config.Config.DbName)
	// log.Println(config.Config.LogFile)
	// fmt.Println(models.Db)

	u := models.User{}
	u.Name = "testName"
	u.Email = "test@example.com"
	u.Password = "testPassword"
	u.CreateUser()
	log.Println("ユーザー作成完了")

	// 作成したユーザーを取得
	user, _ := models.GetUser(1)
	log.Println("取得したユーザー:", user)

	// ユーザー情報を更新
	user.Name = "testName2"
	user.Email = "test2@example.com"
	user.UpdateUser()
	log.Println("ユーザー更新完了")

	// 更新後のユーザーを取得
	user, _ = models.GetUser(1)
	log.Println("更新後のユーザー:", user)
}
