package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func GetConnection() (databse *sql.DB, err error) {
	if db != nil {
		return db, nil
	}

	dpf := mysql.Config{
		User:   "ailre",
		Passwd: "@rtist!nside",
		Net:    "tcp",
		Addr:   "192.168.0.27:3306",
		DBName: "inamazing_test",
	}
	dpf.AllowNativePasswords = true
	spf := dpf.FormatDSN()

	var e error
	db, e = sql.Open("mysql", spf)
	if e != nil {
		log.Fatal(e)
		return nil, e
	}
	return db, nil
}

func DestroyConnection() {
	if db == nil {
		return
	}

	db.Close()
}

// func Join() {
// 	rows, err := db.Query("SELECT id, gtoken, skey from user_cred")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		id, gtoken, skey := new(string), new(string), new(string)
// 		err = rows.Scan(&id, gtoken, skey)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%s\n%s\n%s", *id, *gtoken, *skey)
// 	}
// }
