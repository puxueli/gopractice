package main

import (
	"fmt"
	"net/http"

	"github.com/gohouse/gorose"

	_ "github.com/go-sql-driver/mysql"
)

var err error
var engin *gorose.Engin

func init() {
	engin, err = gorose.Open(&gorose.Config{Driver: "mysql", Dsn: "root:root@(127.0.0.1:3306)/test"})
}

func DB() gorose.IOrm {
	return engin.NewOrm()
}
func main() {
	http.HandleFunc("/", test)
	http.ListenAndServe("127.0.0.1:8000", nil)
}

func test(w http.ResponseWriter, r *http.Request) {
	// orm链式操作,查询单条数据
	res, _ := DB().Table("employinfo").Get()
	// res 类型为 map[string]interface{}
	fmt.Fprintln(w, res)
}
