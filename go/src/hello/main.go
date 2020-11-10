package main

import (
	_ "hello/routers"

	"github.com/astaxie/beego"
)

func main() {
	//	StaticDir["/static"] = "static"'
	beego.SetStaticPath("/static", "static")
	beego.Run()
}
