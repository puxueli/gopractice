// helloworld project main.go
package main

import (
	"fmt"
	"net/http"
)

/*

 */
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
func main() {
	var age int = 20 /*声明实际变量*/
	var address *int /* 声明指针变量 */
	address = &age   /* 指针变量的存储地址 */
	fmt.Println("*age的指针地址是", &age)
	fmt.Println("address的指针是", address)

}

/*
	开启服务
*/
//func main() {
//	http.HandleFunc("/", test)
//	http.ListenAndServe("127.0.0.1:8000", nil)
//}
