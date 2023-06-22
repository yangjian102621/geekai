package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// CopyObject 拷贝对象
func CopyObject(src interface{}, dst interface{}) error {

	srcType := reflect.TypeOf(src)
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst).Elem()
	reflect.TypeOf(dst)
	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		value := dstValue.FieldByName(field.Name)
		if !value.IsValid() {
			continue
		}
		// 数据类型相同，直接赋值
		v := srcValue.FieldByName(field.Name)
		if value.Type() == field.Type {
			value.Set(v)
		} else {
			// src data type is  string，dst data type is slice, map, struct
			// use json decode the data
			if field.Type.Kind() == reflect.String && (value.Type().Kind() == reflect.Struct ||
				value.Type().Kind() == reflect.Map ||
				value.Type().Kind() == reflect.Slice) {
				pType := reflect.New(value.Type())
				v2 := pType.Interface()
				err := json.Unmarshal([]byte(v.String()), &v2)
				if err == nil {
					value.Set(reflect.ValueOf(v2).Elem())
				}
				// map, struct, slice to string
			} else if (field.Type.Kind() == reflect.Struct ||
				field.Type.Kind() == reflect.Map ||
				field.Type.Kind() == reflect.Slice) && value.Type().Kind() == reflect.String {
				ba, err := json.Marshal(v.Interface())
				if err == nil {
					val := string(ba)
					if strings.Contains(val, "{") {
						value.Set(reflect.ValueOf(string(ba)))
					} else {
						value.Set(reflect.ValueOf(""))
					}
				}
			} else { // 简单数据类型的强制类型转换
				switch value.Kind() {
				case reflect.Int:
				case reflect.Int8:
				case reflect.Int16:
				case reflect.Int32:
				case reflect.Int64:
					value.SetInt(v.Int())
					break
				case reflect.Float32:
				case reflect.Float64:
					value.SetFloat(v.Float())
					break
				case reflect.Bool:
					value.SetBool(v.Bool())
					break
				}
			}

		}
	}

	return nil
}

func Ip2Region(searcher *xdb.Searcher, ip string) string {
	str, err := searcher.SearchByStr(ip)
	if err != nil {
		return ""
	}
	arr := strings.Split(str, "|")
	if len(arr) < 3 {
		return arr[0]
	}
	return fmt.Sprintf("%s-%s-%s", arr[0], arr[2], arr[3])
}
