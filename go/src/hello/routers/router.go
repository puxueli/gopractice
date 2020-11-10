package routers

import (
	"hello/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//		beego.Router("/", &controllers.MainController{})
	//		beego.Router("/", &controllers.UserController{}, "*:Get")
	//	beego.Router("/user", &controllers.UserController{}, "*:user")
	beego.Router("user", &controllers.UserController{})
	beego.Router("list", &controllers.ListController{})
	//	beego.Router("/user/list", &UserController{}, "get:list")
}
