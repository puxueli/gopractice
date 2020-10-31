package t

import (
	"fmt"
	"testing"
)

func TestNewT(t *testing.T) {
	var b Type

	//b = New("abc")
	//b = New(1.3)
	//b = New("2.3")
	//b = New(23)
	//b = New(true)
	//b = New(false)
	//b = New(nil)
	//b = New(struct {A string}{"bcd"})
	//b = New(New(3))
	b = New(map[Type]Type{New("a"): New(3)})
	t.Log(b.MapIntT64T()[0])

	t.Log(b.String())
	t.Log(b.Float64())
	t.Log(b.Float32())
	t.Log(b.Int64())
	t.Log(b.Int())
	t.Log(b.Int32())
	t.Log(b.Int16())
	t.Log(b.Int8())
	t.Log(b.Uint64())
	t.Log(b.Uint())
	t.Log(b.Uint32())
	t.Log(b.Uint16())
	t.Log(b.Uint8())
	t.Log(b.Bool())
}

func TestTypeContext_MapStr(t *testing.T) {
	var m = New(map[string]interface{}{"a": 2, "b": "c", "d": "d"})
	a := m.MapStringT()
	t.Log(a)
	t.Log(a["a"])
}

func TestTypeContext_Map(t *testing.T) {
	var m = New(map[interface{}]interface{}{"a": 2, 2: "c", 3.3: "d", true: true, false: 3})
	a := m.MapInterfaceT()
	t.Log(a)
	t.Log(a["a"])
	t.Log(a[2].String())
	t.Log(a[3.3])
	t.Log(a[true].Bool())
	t.Log(a[false])
}

func TestTypeContext_Map2(t *testing.T) {
	var m = New(`{"a": 2, "b":3,"33":{"331":"d"}}`)
	a := m.Map()
	t.Log(a[New("a")])
	t.Log(m.Extract("33.331"))
}

func TestTypeContext_Slice(t *testing.T) {
	var a = New([]string{"a", "b"})

	for _, v := range a.Slice() {
		t.Log(v.String(), v.Bool())
	}
}

func TestTypeContext_Bind(t *testing.T) {
	type json struct {
		A interface{} `json:"a"`
		B string      `json:"b"`
	}
	var a = map[string]interface{}{"a": 1, "b": "bbb"}
	var js json
	New(a).Bind(&js)
	t.Logf("%+v", js)
}

func BenchmarkTypeContext_Int64(b *testing.B) {
	var a Type = New("2.3")
	for i := 0; i < b.N; i++ {
		a.Int64()
	}
}

func TestTypeContext_Map3(t *testing.T) {
	//res := New([]string{"a","b"})
	res := New(`{"aa":11,"bb":["a","b"]}`)
	//res := New([]map[string]interface{}{
	//	{"aa":11,"bb":[]string{"a","b"}},
	//})
	r := res.MapStringInterface()
	fmt.Println(r)
	fmt.Println("res.IsJsonSlice:", res.IsJsonMap())
	fmt.Printf("res.Extract-0.bb.0 : %#v\n", res.Extract("0.bb.0").String())
	fmt.Printf("res.Extract -bb: %#v\n", res.Extract("bb.0").String())
}
