package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type UserController struct {
	beego.Controller
}

type ListController struct {
	beego.Controller
}

//数据表字段定义
type Employinfo struct {
	Id    int64
	Uname string
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1)/test?charset=utf8")
	orm.RegisterModel(new(Employinfo))
	_ = orm.RunSyncdb("default", false, true)
}

func (this *UserController) Get() {
	o := orm.NewOrm()
	user := Employinfo{Id: 5}
	err := o.Read(&user)

	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	} else {
		fmt.Println(user.Id, user.Uname)
	}
	this.Data["id"] = user.Id
	this.Data["name"] = user.Uname
	this.TplName = "user.tpl"

}

func (this *ListController) Get() {
	o := orm.NewOrm()
	var students = []*Employinfo{}
	//	var data map[]
	_, res := o.QueryTable("employinfo").All(&students)

	//接口返回数据
	//	if res == nil {
	//		fmt.Println(students)
	//		jsons, _ := json.Marshal(students)
	//		fmt.Println(string(jsons))
	//		this.Data["json"] = map[string]interface{}{
	//			"code": 0,
	//			"msg":  "获取成功",
	//			"data": students,
	//		}
	//		this.ServeJSON()
	//		this.StopRun()
	//	}

	//模板返回数据
	var data string
	if res == nil {
		//		fmt.Println(students)
		jsons, _ := json.Marshal(students)
		data = string(jsons)

		this.Data["res"] = data
		//return render(request,'list.tpl',{'user_list':data})
		this.TplName = "list.tpl"

	}

}

func user() {
	fmt.Println("313213123123")
}
