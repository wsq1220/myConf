package conf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

// 解析日志库的配置文件

type Config struct {
	FilePath string `conf:"file_path" db:"name"`
	FileName string `conf:"file_name"`
	MaxSize  int64  `conf:"max_size"`
}

// 从conf文件中读取内容赋值给结构体指针
func ParseConf(fileName string, result interface{}) (err error) {
	// 必须是pointer
	t := reflect.TypeOf(result)
	v := reflect.ValueOf(result)

	if t.Kind() != reflect.Ptr {
		err = errors.New("result必须是一个指针")
		return
	}
	// struct
	tElem := t.Elem()
	if tElem.Kind() != reflect.Struct {
		err = errors.New("result必须是一个结构体指针")
		return
	}

	data, err := ioutil.ReadFile(fileName) // []byte
	if err != nil {
		err = fmt.Errorf("打开配置文件%s失败", fileName)
		return err
	}

	lineSlice := strings.Split(string(data), "\r\n")
	fmt.Println(lineSlice)
	// 一行一行的解析
	for index, line := range lineSlice {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			// 忽略空行和注释
			continue
		}

		equalIndex := strings.Index(line, "=")
		if equalIndex == -1 {
			err = fmt.Errorf("第%d行语法错误", index+1)
			return
		}

		key := line[:equalIndex]
		value := line[equalIndex+1:]
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(key) == 0 {
			err = fmt.Errorf("第%d行语法错误", index+1)
			return
		}

		for i := 0; i < tElem.NumField(); i++ {
			field := tElem.Field(i)
			tag := field.Tag.Get("conf")
			if key == tag {

				fieldType := field.Type
				switch fieldType.Kind() {
				case reflect.String:
					fieldValue := v.Elem().FieldByName(field.Name)
					fieldValue.SetString(value)
				case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
					value64, _ := strconv.ParseInt(value, 10, 64)
					v.Elem().Field(i).SetInt(value64)
				}
			}
		}

	}
	return
}
