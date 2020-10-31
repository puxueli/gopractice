package main

import (
	"fmt"
	"net/http"

	"github.com/gohouse/gorose"

	_ "github.com/go-sql-driver/mysql"
)

var err error
var engin *gorose.Engin

type employinfo struct {
	uname   string
	age     int
	phone   string
	address string
}

func init() {

	engin, err = gorose.Open(&gorose.Config{Driver: "mysql", Dsn: "root:root@(127.0.0.1:3306)/test"})
}

func DB() gorose.IOrm {
	return engin.NewOrm()
}
func main() {
	http.HandleFunc("/", del)
	http.ListenAndServe("127.0.0.1:8000", nil)
}

func sel(w http.ResponseWriter, r *http.Request) {
	// orm链式操作,查询单条数据
	res, _ := DB().Table("employinfo").Get()
	fmt.Fprintln(w, res)
}

func add(w http.ResponseWriter, r *http.Request) {
	var data = map[string]interface{}{"uname": "gotestname", "age": 120, "phone": "13312345678", "address": "天上人间"}
	//	fmt.Fprintln(w, "插入数据是", employ)
	id, err := DB().Table("employinfo").Data(data).Insert()
	if err != nil {
		fmt.Fprintln(w, "DB.Insert err:", err)
	}
	fmt.Fprintln(w, "插入了数据，id是：", id)
}

func del(w http.ResponseWriter, r *http.Request) {
	//	var where map[string]interface{}{"id": 27}
	res, err := DB().Table("employinfo").Where("id", 30).Delete()
	if err != nil {
		fmt.Fprintln(w, "DB.del err:", err)
	}
	fmt.Fprintln(w, "删除了数据", res)
}
