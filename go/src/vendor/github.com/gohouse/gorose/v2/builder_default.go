package gorose

import (
	"errors"
	"fmt"
	"github.com/gohouse/gocar/structEngin"
	"github.com/gohouse/t"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var operator = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like",
	"in", "not in", "between", "not between"}

type BuilderDefault struct {
	IOrm
	operator    []string
	placeholder int
	driver      string
}

var onceBuilderDefault sync.Once
var builderDefault *BuilderDefault

func NewBuilderDefault(o IOrm) *BuilderDefault {
	onceBuilderDefault.Do(func() {
		builderDefault = new(BuilderDefault)
		builderDefault.operator = operator
		builderDefault.driver = "mysql"
	})
	builderDefault.IOrm = o
	// 每次使用的时候, 重置为0, 方便pg的占位符使用
	builderDefault.placeholder = 0
	return builderDefault
}

// SetDriver 设置驱动, 方便获取占位符使用
func (b *BuilderDefault) SetDriver(dr string) *BuilderDefault {
	b.driver = dr
	return b
}

// GetPlaceholder 获取占位符
func (b *BuilderDefault) GetPlaceholder() (phstr string) {
	switch b.driver {
	case "postgres":
		withLockContext(func() {
			ph := b.placeholder + 1
			phstr = fmt.Sprintf("$%v", ph)
			b.placeholder = ph
		})
	default:
		phstr = "?"
	}
	return
}

func (b *BuilderDefault) BuildQuery() (sqlStr string, args []interface{}, err error) {
	//b.IOrm = o
	join, err := b.BuildJoin()
	if err != nil {
		b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
		return
	}
	where, err := b.BuildWhere()
	if err != nil {
		b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
		return
	}
	sqlStr = fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
		b.BuildDistinct(), b.BuildFields(), b.BuildTable(), join, where,
		b.BuildGroup(), b.BuildHaving(), b.BuildOrder(), b.BuildLimit(), b.BuildOffset())

	//args = b.bindParams
	args = b.IOrm.GetBindValues()
	return
}

// BuildExecut : build execute query string
func (b *BuilderDefault) BuildExecute(operType string) (sqlStr string, args []interface{}, err error) {
	// insert : {"name":"fizz, "website":"fizzday.net"} or {{"name":"fizz2", "website":"www.fizzday.net"}, {"name":"fizz", "website":"fizzday.net"}}}
	// update : {"name":"fizz", "website":"fizzday.net"}
	// delete : ...
	//b.IOrm = o
	var update, insertkey, insertval string
	if operType != "delete" {
		if b.IOrm.GetData() == nil {
			err = errors.New("insert,update请传入数据操作")
			b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
			return
		}
		update, insertkey, insertval = b.BuildData(operType)
	}

	where, err := b.BuildWhere()
	if err != nil {
		b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
		return
	}

	switch operType {
	case "insert":
		sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", b.BuildTable(), insertkey, insertval)
	case "update":
		if where == "" && b.IOrm.GetForce() == false {
			err = errors.New("出于安全考虑, update时where条件不能为空, 如果真的不需要where条件, 请使用Force()(如: db.xxx.Force().Update())")
			b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
			return
		}
		sqlStr = fmt.Sprintf("UPDATE %s SET %s%s", b.BuildTable(), update, where)
	case "delete":
		if where == "" && b.IOrm.GetForce() == false {
			err = errors.New("出于安全考虑, delete时where条件不能为空, 如果真的不需要where条件, 请使用Force()(如: db.xxx.Force().Delete())")
			b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
			return
		}
		sqlStr = fmt.Sprintf("DELETE FROM %s%s", b.BuildTable(), where)
	}

	args = b.IOrm.GetBindValues()
	return
}

// BuildData : build inert or update data
func (b *BuilderDefault) BuildData(operType string) (string, string, string) {
	data := b.IOrm.GetData()
	ref := reflect.Indirect(reflect.ValueOf(data))

	switch ref.Kind() {
	case reflect.Struct:
		return b.parseData(operType, structEngin.New().SetExtraCols(b.IOrm.GetExtraCols()).StructContent2Map(data))
	case reflect.Map:
		var tmp = []map[string]interface{}{t.New(data).MapStringInterface()}
		return b.parseData(operType, tmp)
	case reflect.Slice:
		switch ref.Type().Elem().Kind() {
		case reflect.Struct:
			return b.parseData(operType, structEngin.New().SetExtraCols(b.IOrm.GetExtraCols()).StructContent2Map(data))
		case reflect.Map:
			return b.parseData(operType, t.New(data).SliceMapStringInterface())
		}
	}
	return "", "", ""
}
func (b *BuilderDefault) BuildData2(operType string) (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	data := b.IOrm.GetData()

	switch data.(type) {
	case string:
		dataObj = append(dataObj, data.(string))
	case []map[string]interface{}, []Data: // insert multi datas ([]map[string]interface{})
		sliceData := t.New(data).Slice()
		for key, _ := range sliceData[0].MapString() {
			if inArray(key, dataFields) == false {
				dataFields = append(dataFields, key)
			}
		}
		for _, itemT := range sliceData {
			item := itemT.MapString()
			var dataValuesSub []string
			for _, key := range dataFields {
				if item[key] == nil {
					dataValuesSub = append(dataValuesSub, "null")
				} else {
					dataValuesSub = append(dataValuesSub, b.GetPlaceholder())
					b.IOrm.SetBindValues(item[key])
				}
			}
			dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
		}
	case map[string]interface{}, Data: // update or insert (map[string]interface{})
		var dataValuesSub []string
		for key, val := range t.New(data).MapString() {
			if operType == "insert" {
				// insert
				dataFields = append(dataFields, key)
				if val.Interface() == nil {
					dataValuesSub = append(dataValuesSub, "null")
				} else {
					dataValuesSub = append(dataValuesSub, b.GetPlaceholder())
					b.IOrm.SetBindValues(val.Interface())
				}
			} else if operType == "update" {
				// update
				if val.Interface() == nil {
					dataObj = append(dataObj, key+"=null")
				} else {
					dataObj = append(dataObj, key+"="+b.GetPlaceholder())
					b.IOrm.SetBindValues(val.Interface())
				}
			}
		}
		if operType == "insert" {
			// insert
			dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
		}
	default:
		//ref := reflect.Indirect(reflect.ValueOf(data))
		//switch ref.Kind() {
		//case reflect.Struct:
		//	structEngin.New().StructContent2Map(data)
		//case reflect.Map:
		//case reflect.Slice:
		//	switch ref.Type().Elem().Kind() {
		//	case reflect.Struct:
		//	case reflect.Map:
		//	}
		//}

		return "", "", ""
	}

	return strings.Join(dataObj, ","), strings.Join(dataFields, ","), strings.Join(dataValues, ",")
}

func (b *BuilderDefault) parseData(operType string, data []map[string]interface{}) (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	for key, _ := range data[0] {
		if inArray(key, dataFields) == false {
			dataFields = append(dataFields, key)
		}
	}
	for _, item := range data {
		// 定义1条数据的存储
		var dataValuesSub []string
		for _, key := range dataFields {
			if item[key] == nil {
				// 放入占位符
				dataValuesSub = append(dataValuesSub, b.GetPlaceholder())
				// 保存真正的值为null
				b.IOrm.SetBindValues("null")
			} else {
				// 放入占位符
				dataValuesSub = append(dataValuesSub, b.GetPlaceholder())
				// 保存真正的值
				b.IOrm.SetBindValues(item[key])
			}
			// update
			dataObj = append(dataObj, fmt.Sprintf("%s=%s", key, b.GetPlaceholder()))
		}
		dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
	}
	return strings.Join(dataObj, ","), strings.Join(dataFields, ","), strings.Join(dataValues, ",")
}

func (b *BuilderDefault) BuildJoin() (s string, err error) {
	// 用户传入的join参数+join类型
	var join []interface{}
	var returnJoinArr []string
	joinArr := b.GetJoin()

	for _, join = range joinArr {
		var w string
		var ok bool
		// 用户传入 join 的where值, 即第二个参数
		var args []interface{}

		if len(join) != 2 {
			err = errors.New("join conditions are wrong")
			b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
			return
		}

		// 获取真正的用户传入的join参数
		if args, ok = join[1].([]interface{}); !ok {
			err = errors.New("join conditions are wrong")
			b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
			return
		}

		argsLength := len(args)
		switch argsLength {
		case 1: // join字符串 raw
			w = args[0].(string)
		case 2: // join表 + 字符串
			w = args[0].(string) + " ON " + args[1].(string)
		case 4: // join表 + (a字段+关系+a字段)
			w = args[0].(string) + " ON " + args[1].(string) + " " + args[2].(string) + " " + args[3].(string)
		default:
			err = errors.New("join format error")
			b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
			return
		}

		returnJoinArr = append(returnJoinArr, " "+join[0].(string)+" JOIN "+w)
	}

	return strings.Join(returnJoinArr, " "), nil
}

func (b *BuilderDefault) BuildWhere() (where string, err error) {
	var beforeParseWhere = b.IOrm.GetWhere()
	where, err = b.parseWhere(b.IOrm)
	b.IOrm.SetWhere(beforeParseWhere)
	return If(where == "", "", " WHERE "+where).(string), err
}

func (b *BuilderDefault) BuildDistinct() (dis string) {
	return If(b.IOrm.GetDistinct(), "DISTINCT ", "").(string)
}

func (b *BuilderDefault) BuildFields() string {
	if len(b.IOrm.GetFields()) == 0 {
		return "*"
	}
	return strings.Join(b.IOrm.GetFields(), ",")
}

func (b *BuilderDefault) BuildTable() string {
	return b.IOrm.GetTable()
}

func (b *BuilderDefault) BuildGroup() string {
	return If(b.IOrm.GetGroup() == "", "", " GROUP BY "+b.IOrm.GetGroup()).(string)
}

func (b *BuilderDefault) BuildHaving() string {
	return If(b.IOrm.GetHaving() == "", "", " HAVING "+b.IOrm.GetHaving()).(string)
}

func (b *BuilderDefault) BuildOrder() string {
	return If(b.IOrm.GetOrder() == "", "", " ORDER BY "+b.IOrm.GetOrder()).(string)
}

func (b *BuilderDefault) BuildLimit() string {
	if b.IOrm.GetUnion()!=nil {
		return ""
	}
	return If(b.IOrm.GetLimit() == 0, "", " LIMIT "+strconv.Itoa(b.IOrm.GetLimit())).(string)
}

func (b *BuilderDefault) BuildOffset() string {
	if b.BuildLimit() == "" {
		return ""
	}
	if b.IOrm.GetUnion()!=nil {
		return ""
	}
	return If(b.IOrm.GetOffset() == 0, "", " OFFSET "+strconv.Itoa(b.IOrm.GetOffset())).(string)
}

// parseWhere : parse where condition
func (b *BuilderDefault) parseWhere(ormApi IOrm) (string, error) {
	// 取出所有where
	wheres := ormApi.GetWhere()
	// where解析后存放每一项的容器
	var where []string

	for _, args := range wheres {
		// and或者or条件
		var condition = args[0].(string)
		// 统计当前数组中有多少个参数
		params := args[1].([]interface{})
		paramsLength := len(params)

		switch paramsLength {
		case 3: // 常规3个参数:  {"id",">",1}
			res, err := b.parseParams(params, ormApi)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)

		case 2: // 常规2个参数:  {"id",1}
			res, err := b.parseParams(params, ormApi)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)
		case 1: // 二维数组或字符串
			switch paramReal := params[0].(type) {
			case string:
				where = append(where, condition+" ("+paramReal+")")
			case map[string]interface{}: // 一维数组
				var whereArr []string
				for key, val := range paramReal {
					whereArr = append(whereArr, key+"="+b.GetPlaceholder())
					b.IOrm.SetBindValues(val)
				}
				where = append(where, condition+" ("+strings.Join(whereArr, " and ")+")")
			case [][]interface{}: // 二维数组
				var whereMore []string
				for _, arr := range paramReal { // {{"a", 1}, {"id", ">", 1}}
					whereMoreLength := len(arr)
					switch whereMoreLength {
					case 3:
						res, err := b.parseParams(arr, ormApi)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					case 2:
						res, err := b.parseParams(arr, ormApi)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					default:
						return "", errors.New("where data format is wrong")
					}
				}
				where = append(where, condition+" ("+strings.Join(whereMore, " and ")+")")
			case func():
				// 清空where,给嵌套的where让路,复用这个节点
				ormApi.SetWhere([][]interface{}{})

				// 执行嵌套where放入Database struct
				paramReal()
				// 再解析一遍后来嵌套进去的where
				wherenested, err := b.parseWhere(ormApi)
				if err != nil {
					b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
					return "", err
				}
				// 嵌套的where放入一个括号内
				where = append(where, condition+" ("+wherenested+")")
			default:
				return "", errors.New("where data format is wrong")
			}
		}
	}

	// 合并where,去掉左侧的空格,and,or并返回
	return strings.TrimLeft(
		strings.TrimLeft(
			strings.TrimLeft(
				strings.Trim(
					strings.Join(where, " "),
					" "),
				"and"),
			"or"),
		" "), nil
}

/**
 * 将where条件中的参数转换为where条件字符串
 * example: {"id",">",1}, {"age", 18}
 */
// parseParams : 将where条件中的参数转换为where条件字符串
func (b *BuilderDefault) parseParams(args []interface{}, ormApi IOrm) (s string, err error) {
	paramsLength := len(args)
	argsReal := args

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		//if !inArray(argsReal[1], b.GetRegex()) {
		if !inArray(argsReal[1], b.GetOperator()) {
			err = errors.New("where parameter is wrong")
			b.IOrm.GetISession().GetIEngin().GetLogger().Error(err.Error())
			return
		}

		paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, argsReal[1].(string))

		switch strings.Trim(strings.ToLower(t.New(argsReal[1]).String())," ") {
		//case "like", "not like":
		//	paramsToArr = append(paramsToArr, b.GetPlaceholder())
		//	b.IOrm.SetBindValues(argsReal[2])
		case "in", "not in":
			var tmp []string
			var ar2 = t.New(argsReal[2]).Slice()
			for _, item := range ar2 {
				tmp = append(tmp, b.GetPlaceholder())
				b.IOrm.SetBindValues(t.New(item).Interface())
			}
			paramsToArr = append(paramsToArr, "("+strings.Join(tmp, ",")+")")

		case "between", "not between":
			var ar2 = t.New(argsReal[2]).Slice()
			paramsToArr = append(paramsToArr, b.GetPlaceholder()+" and "+b.GetPlaceholder())
			b.IOrm.SetBindValues(ar2[0].Interface())
			b.IOrm.SetBindValues(ar2[1].Interface())

		default:
			paramsToArr = append(paramsToArr, b.GetPlaceholder())
			b.IOrm.SetBindValues(argsReal[2])
		}
	case 2:
		paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, "=")
		paramsToArr = append(paramsToArr, b.GetPlaceholder())
		b.IOrm.SetBindValues(argsReal[1])
	}

	return strings.Join(paramsToArr, " "), nil
}

func (b *BuilderDefault) GetOperator() []string {
	return b.operator
}
