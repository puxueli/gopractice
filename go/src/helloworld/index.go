package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", test)
	http.ListenAndServe("127.0.0.1:8000", nil)

}

func test(w http.ResponseWriter, r *http.Request) {
	db, dbconn := sql.Open("mysql", "root:root@(127.0.0.1:3306)/test")
	fmt.Println(dbconn)
	if dbconn != nil {
		fmt.Println("query error")
		os.Exit(1)
	}
	rows, dbconn := db.Query("select * from employinfo")

	var id int
	var uname string
	var age int
	var phone string
	var address string
	for rows.Next() {

		dbconn = rows.Scan(&id, &uname, &age, &phone, &address)
		fmt.Println(id)
		fmt.Println(uname)
		fmt.Println(age)
		fmt.Println(phone)
	}
}
