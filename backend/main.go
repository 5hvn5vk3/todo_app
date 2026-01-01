package main

import (
	"backend/controllers"
	"backend/models"
	"fmt"
)

func main() {
	fmt.Println(models.Db)

	controllers.StartMainServer()
}
