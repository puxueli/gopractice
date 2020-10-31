package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", del)
	http.ListenAndServe("127.0.0.1:7000", nil)

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

func insert(w http.ResponseWriter, r *http.Request) {
	db, dbconn := sql.Open("mysql", "root:root@(127.0.0.1:3306)/test")
	//	fmt.Println(dbconn)
	if dbconn != nil {
		fmt.Println("query error")
		os.Exit(1)
	}
	result, err := db.Exec("INSERT INTO employinfo(uname,age,phone,address)VALUES(?,?,?,?)", "gotest", "100", "13212345678", "北京市朝阳区")
	if err != nil {
		fmt.Println("insert failed,", err)
	}
	userId, err := result.LastInsertId()   //获取添加数据的id
	rowCount, err := result.RowsAffected() //影响行数
	fmt.Println("user_id:", userId)
	fmt.Println("rowCount:", rowCount)
}

func del(w http.ResponseWriter, r *http.Request) {
	db, dbconn := sql.Open("mysql", "root:root@(127.0.0.1:3306)/test")
	//	fmt.Println(dbconn)
	if dbconn != nil {
		fmt.Println("query error")
		os.Exit(1)
	}

	result, err := db.Exec("delete from employinfo where id=23")
	if err != nil {
		fmt.Println("insert failed,", err)
	}
	rowCount, err := result.RowsAffected() //影响行数
	fmt.Println("rowCount:", rowCount)
}
