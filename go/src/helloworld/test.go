package main

import "fmt"

///*
//	定义结构体
//*/
//type Books struct {
//	title string
//	//	author string
//	//	subject string
//	//	book_id int
//}

//func test() {

//	var book Books //实例化结构体类型
//	var book1 Books
//	book.title = "GOYUYAN"
//	book1.title = "PGP"
//	fmt.Println(book.title)
//	fmt.Println(book1.title)
//}

//func qiepian() {
//	s := []int{1, 2, 3}
//	fmt.Println(s)
//}
//func main() {
//	qiepian()
//}

func fibonacci(n int) int {
	if n < 3 {
		return n
	}
	return fibonacci(n-2) + fibonacci(n-1)
}

func main() {
	var i int
	for i = 1; i <= 5; i++ {
		fmt.Printf("%d\t", fibonacci(i))
	}
}
