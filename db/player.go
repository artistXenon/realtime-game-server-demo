package db

import (
	"log"
)

func GetPlayer(id string, lobby string) (session *string) {
	db, err := GetConnection()
	if err != nil {
		return nil
	}
	var session_key string
	err = db.QueryRow("SELECT skey FROM user_cred WHERE id=? AND lid=?", id, lobby).Scan(&session_key)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &session_key
}
