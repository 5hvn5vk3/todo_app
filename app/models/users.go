package models

import (
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (u *User) CreateUser() (err error) {
	cmd := `insert into users(
	    uuid,
		name,
		email,
		password,
		created_at) values(?,?,?,?,?)`
	Db.Exec(cmd,
		createUUID,
		u.Name,
		u.Email,
		Encrypt(u.Password),
		time.Now())
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
