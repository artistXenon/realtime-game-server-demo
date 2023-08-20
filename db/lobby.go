package db

import (
	"log"
)

func DeleteLobby(id string) {
	db, err := GetConnection()
	if err != nil {
		return
	}
	log.Printf("delete performed %s\n", id)
	_, err = db.Exec("DELETE FROM lobby WHERE id=?", id)
	if err != nil {
		log.Print(err)
		return
	}
	return
}
