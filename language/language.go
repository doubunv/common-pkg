package language

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"reflect"
	"runtime/debug"
	"strings"
)

// 全局定义,会在IP解析的时候使用
var LanguageStandMap = map[string]string{
	"巴西": "pt-BR",
}

const EnglishLanguage = "en-US"
const ChineseLanguage = "zh-CN"
const IndiaLanguage = "hi-IN"
const BRLanguage = "pt-BR"

var indiaLanguageMap map[string]interface{}
var chineseLanguageMap map[string]interface{}
var bRLanguageMap map[string]interface{}

func init() {
	json.Unmarshal([]byte(chineseLanguage), &chineseLanguageMap)
}

func SwitchLanguageDic(data string, language string) interface{} {
	if language == EnglishLanguage {
		return data
	}
	languageData := getLanguageData(language)
	if newValue, ok := languageData[data]; ok {
		return newValue
	}
	return data
}

func SwitchLanguage(data interface{}, language string) interface{} {
	defer func() {
		if err := recover(); err != nil {
			logc.Errorf(context.Background(), "ConsumeMessagesWithContext error:%v, %s", err, string(debug.Stack()))
			return
		}
	}()

	if language == EnglishLanguage {
		return data
	}
	return recursiveGetAllValues(data, getLanguageData(language))
}

func getLanguageData(language string) map[string]interface{} {
	switch language {
	case IndiaLanguage:
		return indiaLanguageMap
	case ChineseLanguage:
		return chineseLanguageMap
	case BRLanguage:
		return bRLanguageMap
	}
	return nil
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
