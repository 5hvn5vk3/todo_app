package main

import (
	"fmt"
	"todo_app/app/models"
	_ "todo_app/config"
)

func main() {
	// log.Println(config.Config.Port)
	// log.Println(config.Config.SQLDriver)
	// log.Println(config.Config.DbName)
	// log.Println(config.Config.LogFile)
	// fmt.Println(models.Db)

	// u := models.User{}
	// u.Name = "testName"
	// u.Email = "test@example.com"
	// u.Password = "testPassword"
	// log.Println(u)

	// u.CreateUser()

	u, _ := models.GetUser(1)
	fmt.Println(u)
}
