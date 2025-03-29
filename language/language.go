package language

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
)

var languageMap map[string]map[string]interface{}

func init() {
	dirPath := "etc/language"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("读取文件夹时发生错误:", err)
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			readJson(dirPath + "/" + file.Name())
		}
	}
}

func readJson(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("文件 %s 不存在\n", filePath)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("打开文件时发生错误:", err)
		return
	}
	defer file.Close() // 确保文件最终会关闭
	var data map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("解析 JSON 时发生错误:", err)
		return
	}

	baseName := filepath.Base(file.Name())
	fileNameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
	languageMap[fileNameWithoutExt] = data
}

func SwitchLanguage(data interface{}, language string) interface{} {
	defer func() {
		if err := recover(); err != nil {
			logc.Errorf(context.Background(), "ConsumeMessagesWithContext error:%v, %s", err, string(debug.Stack()))
			return
		}
	}()

	if _, ok := languageMap[language]; !ok {
		return data
	}

	return recursiveGetAllValues(data, languageMap[language])
}

func checkType(data interface{}) bool {
	sw := reflect.ValueOf(data).Kind()
	if sw == reflect.String {
		return true
	} else {
		return false
	}
}

func recursiveGetAllValues(data interface{}, newData map[string]interface{}) interface{} {

	if newData == nil {
		return data
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		result := make(map[string]interface{})
		keys := v.MapKeys()
		for _, key := range keys {
			value := recursiveGetAllValues(v.MapIndex(key).Interface(), newData)
			result[key.String()] = value
			if checkType(v.MapIndex(key).Interface()) {
				var bl = false
				if strings.ToLower(key.String()) == "code" ||
					strings.ToLower(key.String()) == "id" {
					bl = true
				}

				if newValue, ok := newData[value.(string)]; ok && !bl {
					result[key.String()] = newValue
				}
			}
		}
		return result

	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			value := recursiveGetAllValues(v.Index(i).Interface(), newData)
			result[i] = value
			if checkType(v.Index(i).Interface()) {
				if newValue, ok := newData[value.(string)]; ok {
					result[i] = newValue
				}
			}
		}
		return result

	case reflect.Struct:
		result := make(map[string]interface{})
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			value := recursiveGetAllValues(v.Field(i).Interface(), newData)

			fieldTag := strings.Split(field.Tag.Get("json"), ",")[0]
			result[fieldTag] = value
			if checkType(v.Field(i).Interface()) {
				var bl = false
				if strings.ToLower(fieldTag) == "code" ||
					strings.ToLower(fieldTag) == "id" {
					bl = true
				}
				if newValue, ok := newData[value.(string)]; ok && !bl {
					result[fieldTag] = newValue
				}
			}
		}
		return result

	case reflect.Interface:
		if !v.IsNil() {
			return recursiveGetAllValues(v.Elem().Interface(), newData)
		}
	}

	return data
}
