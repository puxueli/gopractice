package main

import (
	"fmt"
)

func main() {
	n := 5
	js(n)
}

func js(n int) {
	var i int
	var res int = 0
	var resut int = 0
	for i = 1; i <= 5; i++ {
		if i >= 3 {
			res = (i - 1) + (i - 2)
		} else {
			res = 1
		}
		resut = resut + res
	}

	fmt.Println(resut)

}
