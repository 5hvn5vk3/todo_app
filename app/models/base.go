package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"os"
	"todo_app/config"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

var Db *sql.DB

var err error

func init() {

	url := os.Getenv("DATABASE_URL")
	connection, _ := pq.ParseURL(url)
	connection += "sslmode=require"
	Db, err = sql.Open(config.Config.SQLDriver, connection)
	if err != nil {
		log.Fatalln(err)
	}
}

func createUUID() (uuidobj uuid.UUID) { // update: users.go から移動
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(plaintext string) (cryptext string) { // update: users.go から移動
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
