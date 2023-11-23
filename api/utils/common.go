package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"reflect"
	"strconv"
	"strings"
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
				if err == nil && v2 != nil {
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
			} else if field.Type.Kind() != value.Type().Kind() { // 不同类型的字段过滤掉
				continue
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

func IsEmptyValue(obj interface{}) bool {
	if obj == nil {
		return true
	}

	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	default:
		return reflect.DeepEqual(obj, reflect.Zero(reflect.TypeOf(obj)).Interface())
	}
}

func BoolValue(str string) bool {
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return value
}

func FloatValue(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return value
}

func IntValue(str string, defaultValue int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return value
}

func ForceCovert(src any, dst interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, dst)
	if err != nil {
		return err
	}
	return nil
}

func GenQrcode(text string, size int, logo io.Reader) ([]byte, error) {
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	qr.BackgroundColor = color.White
	qr.ForegroundColor = color.Black
	if logo == nil {
		return qr.PNG(size)
	}

	// 生成带Logo的二维码图像
	logoImage, _, err := image.Decode(logo)
	if err != nil {
		return nil, err
	}

	// 缩放 Logo
	scaledLogo := resize.Resize(uint(size/9), uint(size/9), logoImage, resize.Lanczos3)
	// 将Logo叠加到二维码图像上
	qrWithLogo := overlayLogo(qr.Image(size), scaledLogo)

	// 将带Logo的二维码图像以JPEG格式编码为图片数据
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, qrWithLogo, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 叠加Logo到图片上
func overlayLogo(qrImage, logoImage image.Image) image.Image {
	offsetX := (qrImage.Bounds().Dx() - logoImage.Bounds().Dx()) / 2
	offsetY := (qrImage.Bounds().Dy() - logoImage.Bounds().Dy()) / 2

	combinedImage := image.NewRGBA(qrImage.Bounds())
	draw.Draw(combinedImage, qrImage.Bounds(), qrImage, image.Point{}, draw.Over)
	draw.Draw(combinedImage, logoImage.Bounds().Add(image.Pt(offsetX, offsetY)), logoImage, image.Point{}, draw.Over)

	return combinedImage
}
