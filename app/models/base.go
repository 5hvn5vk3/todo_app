package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
)

var Db *sql.DB

var err error

func init() {
	dsn := os.Getenv("DATABASE_URL")
	Db, err = sql.Open("postgres", dsn)
}

func createUUID() (uuidobj uuid.UUID) { // update: users.go から移動
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(plaintext string) (cryptext string) { // update: users.go から移動
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
