# GoRose ORM
[![GoDoc](https://godoc.org/github.com/gohouse/gorose?status.svg)](https://godoc.org/github.com/gohouse/gorose)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohouse/gorose)](https://goreportcard.com/report/github.com/gohouse/gorose)
[![GitHub release](https://img.shields.io/github/release/gohouse/gorose.svg)](https://github.com/gohouse/gorose/releases/latest)
[![Gitter](https://badges.gitter.im/gohouse/gorose.svg)](https://gitter.im/gorose/wechat)
<a target="_blank" href="https://jq.qq.com/?_wv=1027&k=5JJOG9E">
<img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="gorose-orm" title="gorose-orm"></a>

```
  _______   ______   .______        ______        _______. _______ 
 /  _____| /  __  \  |   _  \      /  __  \      /       ||   ____|
|  |  __  |  |  |  | |  |_)  |    |  |  |  |    |   (----`|  |__   
|  | |_ | |  |  |  | |      /     |  |  |  |     \   \    |   __|  
|  |__| | |  `--'  | |  |\  \----.|  `--'  | .----)   |   |  |____ 
 \______|  \______/  | _| `._____| \______/  |_______/    |_______|
```

## 简介
gorose是一个golang orm框架, 借鉴自laravel的eloquent. 
gorose 2.0 采用模块化架构, 通过interface的api通信,严格的上层依赖下层.每一个模块都可以拆卸, 甚至可以自定义为自己喜欢的样子.  
模块关系图如下:  ![gorose-2.0-design](https://i.loli.net/2019/06/19/5d0a1273f12ef86624.jpg)

## 安装
- go.mod
```bash
require github.com/gohouse/gorose v2.0.5
```

- docker
```bash
docker run -it --rm ababy/gorose sh -c "go run main.go"
```
> docker 镜像: [ababy/gorose](https://cloud.docker.com/u/ababy/repository/docker/ababy/gorose), docker镜像包含了gorose所必须的包和运行环境, [查看`Dockerfile`](https://github.com/docker-box/gorose/blob/master/master/golang/Dockerfile)   

## 文档
[最新版2.x文档](https://www.kancloud.cn/fizz/gorose-2/1135835)  
[1.x文档](https://www.kancloud.cn/fizz/gorose/769179)  
[0.x文档](https://gohouse.github.io/gorose/dist/en/index.html)

## api预览
```go
db.Table().Fields().Where().GroupBy().Having().OrderBy().Limit().Select()
db.Table().Data().Insert()
db.Table().Data().Where().Update()
db.Table().Where().Delete()
```

## 简单用法示例
```go
package main
import (
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/mattn/go-sqlite3"
)
var err error
var engin *gorose.Engin
func init() {
    // 全局初始化数据库,并复用
    // 这里的engin需要全局保存,可以用全局变量,也可以用单例
    // 配置&gorose.Config{}是单一数据库配置
    // 如果配置读写分离集群,则使用&gorose.ConfigCluster{}
	engin, err = gorose.Open(&gorose.Config{Driver: "sqlite3", Dsn: "./db.sqlite"})
    // mysql示例, 记得导入mysql驱动 github.com/go-sql-driver/mysql
	// engin, err = gorose.Open(&gorose.Config{Driver: "mysql", Dsn: "root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=true"})
}
func DB() gorose.IOrm {
	return engin.NewOrm()
}
func main() {
    // 单条数据
    res, err := DB().Table("users").First()
    // res 类型为 map[string]interface{}
    fmt.Println(res)
    
    // 多条数据
    res2, _ := DB().Table("users").Get()
    // res2 类型为 []map[string]interface{}
    fmt.Println(res2)
}
```

## 使用建议
gorose提供数据对象绑定(map, struct), 同时支持字符串表名和map数据返回. 提供了很大的灵活性  
建议优先采用数据绑定的方式来完成查询操作, 做到数据源类型可控  
gorose提供了默认的 `gorose.Map` 和 `gorose.Data` 类型, 用来方便初始化绑定和data

## 配置和链接初始化
简单配置
```go
var configSimple = &gorose.Config{
	Driver: "sqlite3", 
	Dsn: "./db.sqlite",
}
```
更多配置, 可以配置集群,甚至可以同时配置不同数据库在一个集群中, 数据库会随机选择集群的数据库来完成对应的读写操作, 其中master是写库, slave是读库, 需要自己做好主从复制, 这里只负责读写
```go
var config = &gorose.ConfigCluster{
	Master:       []&gorose.Config{}{configSimple}
    Slave:        []&gorose.Config{}{configSimple}
    Prefix:       "pre_",
    Driver:       "sqlite3",
}
```
初始化使用
```go
var engin *gorose.Engin
engin, err := Open(config)

if err != nil {
    panic(err.Error())
}
```

## 原生sql操作(增删改查), session的使用
创建用户表 `users`
```sql
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
	 "uid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "name" TEXT NOT NULL,
	 "age" integer NOT NULL
);

INSERT INTO "users" VALUES (1, 'gorose', 18);
INSERT INTO "users" VALUES (2, 'goroom', 18);
INSERT INTO "users" VALUES (3, 'fizzday', 18);
```
定义表struct
```go
type Users struct {
	Uid  int    `gorose:"uid"`
	Name string `gorose:"name"`
	Age  int    `gorose:"age"`
}
// 设置表名, 如果没有设置, 默认使用struct的名字
func (u *Users) TableName() string {
	return "users"
}
```
原生查询操作
```go
// 这里是要绑定的结构体对象
// 如果你没有定义结构体, 则可以直接使用map, map示例
// var u = gorose.Data{}
// var u = gorose.Map{}  这两个都是可以的
var u Users
session := engin.NewSession()
// 这里Bind()是为了存放结果的, 如果你使用的是NewOrm()初始化,则可以直接使用 NewOrm().Table().Query()
err := session.Bind(&u).Query("select * from users where uid=? limit 2", 1)
fmt.Println(err)
fmt.Println(u)
fmt.Println(session.LastSql())
```
原生增删改操作
```go
session.Bind(&u).Query("insert into users(name,age) values(?,?)(?,?)", "gorose",18,"fizzday",19)
session.Bind(&u).Query("update users set name=? where uid=?","gorose",1)
session.Bind(&u).Query("delete from users where uid=?", 1)
```
## 对象关系映射, orm的使用
- 1. 基本链式使用
```go
var u Users
db := engin.NewOrm()
err := db.Table(&u).Fields("name").AddFields("uid","age").Distinct().Where("uid",">",0).OrWhere("age",18).
	Group("age").Having("age>1").OrderBy("uid desc").Limit(10).Offset(1).Select()
```
- 2. 如果不想定义struct, 又想绑定指定类型的map结果, 则可以定义map类型, 如
```go
type user gorose.Map
// 或者 以下的type定义, 都是可以正常解析的
type user2 map[string]interface{}
type users3 []user
type users4 []map[string]string
type users5 []gorose.Map
type users6 []gorose.Data
```
- 2.1 开始使用map绑定
```go
db.Table(&user).Select()
db.Table(&users4).Limit(5).Select()
```
> 注意: 如果使用的不是slice数据结构, 则只能获取到一条数据  

---
这里使用的 gorose.Data , 实际上就是 `map[string]interface{}` 类型.  
而 `gorose.Map`, 实际上是 `t.MapString` 类型, 这里出现了一个 `t` 包, 是一个golang基本数据类型的相互转换包, 请看详细介绍 http://github.com/gohouse/t

- 3. laravel的`First()`,`Get()`, 用来返回结果集  
也就是说, 你甚至可以不用传入各种绑定的struct和map, 直接传入表名, 返回两个参数, 一个是 `[]gorose.Map`结果集, 第二个是`error`,看成简单粗暴  
用法就是把上边的 `Select()` 方法换成 Get,First 即可, 只不过, `Select()` 只返回一个参数  

- 4. orm的增删改查  
```go
db.Table(&user2).Limit(10.Select()
db.Table(&user2).Where("uid", 1).Data(gorose.Data{"name","gorose"}).Update()
db.Table(&user2).Data(gorose.Data{"name","gorose33"}).Insert()
db.Table(&user2).Data([]gorose.Data{{"name","gorose33"},"name","gorose44"}).Insert()
db.Table(&user2).Where("uid", 1).Delete()
```

## 最终sql构造器, builder构造不同数据库的sql
目前支持 mysql, sqlite3, postgres, oracle, mssql, clickhouse等符合 `database/sql` 接口支持的数据库驱动  
这一部分, 用户基本无感知, 分理出来, 主要是为了开发者可以自由添加和修改相关驱动以达到个性化的需求  

## binder, 数据绑定对象  
这一部分也是用户无感知的, 主要是传入的绑定对象解析和数据绑定, 同样是为了开发者个性化定制而独立出来的

## 模块化
gorose2.0 完全模块化, 每一个模块都封装了interface接口api, 模块间调用, 都是通过接口, 上层依赖下层

- 主模块
    - engin  
    gorose 初始化配置模块, 可以全局保存并复用
    - session  
    真正操作数据库底层模块, 所有的操作, 最终都会走到这里来获取或修改数据  
    - orm  
    对象关系映射模块, 所有的orm操作, 都在这里完成  
    - builder  
    构建终极执行的sql模块, 可以构建任何数据库的sql, 但要符合`database/sql`包的接口  
- 子模块
    - driver  
    数据库驱动模块, 被engin和builder依赖, 根据驱动来搞事情  
    - binder  
    结果集绑定模块, 所有的返回结果集都在这里  

以上主模块, 都相对独立, 可以个性化定制和替换, 只要实现相应模块的接口即可.  

## 最佳实践
sql
```sql
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
	 "uid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "name" TEXT NOT NULL,
	 "age" integer NOT NULL
);

INSERT INTO "users" VALUES (1, 'gorose', 18);
INSERT INTO "users" VALUES (2, 'goroom', 18);
INSERT INTO "users" VALUES (3, 'fizzday', 18);
```
实战代码
```go
package main

import (
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/mattn/go-sqlite3"
)

type Users struct {
    Uid int64 `gorose:"uid"`
    Name string `gorose:"name"`
    Age int64 `gorose:"age"`
    Xxx interface{} `gorose:"-"` // 这个字段在orm中会忽略
}

func (u *Users) TableName() string {
	return "users"
}

var err error
var engin *gorose.Engin

func init() {
    // 全局初始化数据库,并复用
    // 这里的engin需要全局保存,可以用全局变量,也可以用单例
    // 配置&gorose.Config{}是单一数据库配置
    // 如果配置读写分离集群,则使用&gorose.ConfigCluster{}
	engin, err = gorose.Open(&gorose.Config{Driver: "sqlite3", Dsn: "./db.sqlite"})
}
func DB() gorose.IOrm {
	return engin.NewOrm()
}
func main() {
	// 这里定义一个变量db, 是为了复用db对象, 可以在最后使用 db.LastSql() 获取最后执行的sql
	// 如果不复用 db, 而是直接使用 DB(), 则会新建一个orm对象, 每一次都是全新的对象
	// 所以复用 db, 一定要在当前会话周期内
	db := DB()
	
	// 查询一条
	var u Users
	// 查询数据并绑定到 user{} 上
	err = db.Table(&u).Fields("uid,name,age").Where("age",">",0).OrderBy("uid desc").Select()
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(u, u.Name)
	fmt.Println(db.LastSql())
	
	// 查询多条
	// 查询数据并绑定到 []Users 上, 这里复用了 db 及上下文条件参数
	// 如果不想复用,则可以使用DB()就会开启全新会话,或者使用db.Reset()
	// db.Reset()只会清除上下文参数干扰,不会更换链接,DB()则会更换链接
	var u2 []Users
	err = db.Limit(10).Offset(1).Select()
	fmt.Println(u2)
	
	// 统计数据
	var count int64
	// 这里reset清除上边查询的参数干扰, 可以统计所有数据, 如果不清楚, 则条件为上边查询的条件
	// 同时, 可以新调用 DB(), 也不会产生干扰
	count,err = db.Reset().Count()
	// 或
	count, err = DB().Table(&u).Count()
	fmt.Println(count, err)
}
```

## 高级用法

- Chunk 数据分片 大量数据批量处理 (累积处理)   

   ` 当需要操作大量数据的时候, 一次性取出再操作, 不太合理, 就可以使用chunk方法  
        chunk的第一个参数是指定一次操作的数据量, 根据业务量, 取100条或者1000条都可以  
        chunk的第二个参数是一个回调方法, 用于书写正常的数据处理逻辑  
        目的是做到, 无感知处理大量数据  
        实现原理是, 每一次操作, 自动记录当前的操作位置, 下一次重复取数据的时候, 从当前位置开始取
        `
	```go
	User := db.Table("users")
	User.Fields("id, name").Where("id",">",2).Chunk(2, func(data []map[string]interface{}) {
	    // for _,item := range data {
	    // 	   fmt.Println(item)
	    // }
	    fmt.Println(data)
	})

	// 打印结果:  
	// map[id:3 name:gorose]
	// map[id:4 name:fizzday]
	// map[id:5 name:fizz3]
	// map[id:6 name:gohouse]
	[map[id:3 name:gorose] map[name:fizzday id:4]]
	[map[id:5 name:fizz3] map[id:6 name:gohouse]]
	```
    
- Loop 数据分片 大量数据批量处理 (从头处理)   

	` 类似 chunk 方法, 实现原理是, 每一次操作, 都是从头开始取数据
	原因: 当我们更改数据时, 更改的结果可能作为where条件会影响我们取数据的结果,所以, 可以使用Loop`
    ```go
	User := db.Table("users")
	User.Fields("id, name").Where("id",">",2).Loop(2, func(data []map[string]interface{}) {
	    // for _,item := range data {
	    // 	   fmt.Println(item)
	    // }
	    // 这里执行update / delete  等操作
	})
	```
    
- 嵌套where  

	```go
	// SELECT  * FROM users  
	//     WHERE  id > 1 
	//         and ( name = 'fizz' 
	//             or ( name = 'fizz2' 
	//                 and ( name = 'fizz3' or website like 'fizzday%')
	//                 )
	//             ) 
	//     and job = 'it' LIMIT 1
	User := db.Table("users")
	User.Where("id", ">", 1).Where(func() {
	        User.Where("name", "fizz").OrWhere(func() {
	            User.Where("name", "fizz2").Where(func() {
	                User.Where("name", "fizz3").OrWhere("website", "like", "fizzday%")
	            })
	        })
	    }).Where("job", "it").First()
	```
