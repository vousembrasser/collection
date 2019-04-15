package collection

import (
	"errors"
	"fmt"
	"reflect"
)

type ObjArray struct{
	AbsArray
	objs reflect.Value // 数组对象，是一个slice
	typ reflect.Type // 数组对象每个元素类型
	ptr reflect.Value // 指向数组对象的指针
}

// 根据对象数组创建
func NewObjArray(objs interface{}) *ObjArray {
	vals := reflect.ValueOf(objs)
	typ := reflect.TypeOf(objs).Elem()
	arr := &ObjArray{
		objs: vals,
		typ: typ,
	}
	arr.AbsArray.Parent = arr
	return arr
}

func (arr *ObjArray) Insert(index int, obj interface{}) IArray {
	if arr.Err() != nil {
		return arr
	}

	ret := arr.objs.Slice(0, index)
	length := arr.objs.Len()
	tail := arr.objs.Slice(index, length)
	ret = reflect.Append(ret, reflect.ValueOf(obj))
	for i := 0; i < tail.Len(); i++ {
		ret = reflect.Append(ret, tail.Index(i))
	}
	arr.objs = ret
	arr.AbsArray.Parent = arr
	return arr
}

func (arr *ObjArray) Index(i int) IMix {
	return NewMix(arr.objs.Index(i).Interface())
}

func (arr *ObjArray) NewEmpty(err ...error) IArray {
	objs := reflect.MakeSlice(arr.objs.Type(), 0, 0)
	ret := &ObjArray{
		objs: objs,
		typ: arr.typ,
	}
	ret.AbsArray.Parent = ret
	if len(err) != 0 {
		ret.SetErr(err[0])
	}
	return ret
}

func (arr *ObjArray) Remove(i int) IArray {
	if arr.Err() != nil {
		return arr
	}

	len := arr.Count()
	if i >= len {
		return arr.SetErr(errors.New("index exceeded"))
	}

	ret := arr.objs.Slice(0, i)
	length := arr.objs.Len()
	tail := arr.objs.Slice(i + 1, length)
	for i := 0; i < tail.Len(); i++ {
		ret = reflect.Append(ret, tail.Index(i))
	}
	arr.objs = ret
	arr.AbsArray.Parent = arr
	return arr
}

func (arr *ObjArray) Count() int {
	return arr.objs.Len()
}

func (arr *ObjArray) DD() {
	ret := fmt.Sprintf("ObjArray(%d):{\n", arr.Count())
	for i:= 0; i< arr.objs.Len(); i++ {
		ret = ret + fmt.Sprintf("\t%d:\t%v\n", i, arr.objs.Index(i))
	}
	ret = ret + "}\n"
	fmt.Print(ret)
}

// Column return some key by column
func (arr *ObjArray) Column(key string) IArray {
	var objs IArray

	field, found := arr.typ.FieldByName(key)
	if !found  {
		panic("ObjArray.Column:key not found")
	}

	switch field.Type.Kind() {
	case reflect.String:
		objs = NewStrArray([]string{})
	case reflect.Int64:
		objs = NewInt64Array([]int64{})
	case reflect.Int:
		objs = NewIntArray([]int{})
	case reflect.Float32:
		objs = NewFloat32Array([]float32{})
	case reflect.Float64:
		objs = NewFloat64Array([]float64{})
	default:
		err := errors.New("ObjArray.Column: not support kind")
		arr.SetErr(err)
		return arr
	}

	for i := 0; i < arr.objs.Len(); i++ {
		v := arr.objs.Index(i).FieldByName(key).Interface()
		objs.Append(v)
	}

	return objs
}

func (arr *ObjArray) KeyBy(key string) (IMap, error) {

	field, found := arr.typ.FieldByName(key)
	if !found  {
		err := errors.New("ObjArray.KeyBy: key not found")
		arr.SetErr(err)
		return nil, err
	}
	m := NewEmptyMap(field.Type, arr.typ)
	for i := 0; i < arr.objs.Len(); i++ {
		v := arr.objs.Index(i).FieldByName(key).Interface()
		m.Set(v, arr.objs.Index(i).Interface())
	}
	return m, nil
}

// 将对象的某个key作为Slice的value，作为slice返回
func (arr *ObjArray) Pluck(key string) IArray {
	//TODO: not implement
	panic(1)
}

// 按照某个字段进行排序
func (arr *ObjArray) SortBy(key string) IArray {
	//TODO: not implement
	panic(1)
}

// 按照某个字段进行排序,倒序
func (arr *ObjArray) SortByDesc(key string) IArray {
	//TODO: not implement
	panic(1)
}