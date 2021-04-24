package utils

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/bitly/go-simplejson"
	"github.com/realjf/goframe/pkg/crypto/md5"
)

func createDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("Failed to create directory '%v': %v", path, err)
			}
		} else {
			return fmt.Errorf("Failed to create directory '%v': %v", path, err)
		}
	}
	return nil
}

func Addslashes(v string) string {
	pos := 0
	buf := make([]byte, len(v)*2)
	for i := 0; i < len(v); i++ {
		c := v[i]
		if c == '\'' || c == '"' || c == '\\' {
			buf[pos] = '\\'
			buf[pos+1] = c
			pos += 2
		} else {
			buf[pos] = c
			pos++
		}
	}
	return string(buf[:pos])
}

func IsValid(name string) bool {
	if match, _ := regexp.MatchString("^[a-z]([-a-z0-9]*[a-z0-9])?$", name); !match {
		return false
	}
	if len(name) <= 2 || len(name) > 20 {
		return false
	}
	return true
}

func IsMysqlName(name string) bool {
	if match, _ := regexp.MatchString("^[a-z]([_a-z0-9]*[a-z0-9])?$", name); !match {
		return false
	}
	if len(name) <= 2 || len(name) > 20 {
		return false
	}
	return true
}

func IsPasswd(name string) bool {
	if match, _ := regexp.MatchString("^[A-Za-z0-9_#!]*$", name); !match {
		return false
	}
	if len(name) <= 5 || len(name) > 20 {
		return false
	}
	return true
}

func IsIP(ip string) (b bool) {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		return false
	}
	return true
}

func JsonToArr(json string) []interface{} {
	var data []interface{}
	sj, err := simplejson.NewJson([]byte(json))
	if err == nil {
		data, _ = sj.Array()
	}
	return data
}

func ToString(arg interface{}) string {
	switch arg.(type) {
	case int, int8, int16:
		return strconv.Itoa(arg.(int))
	case int32, int64:
		return strconv.FormatInt(arg.(int64), 10)
	case []byte:
		return string(arg.([]byte))
	case float32, float64:
		return strconv.FormatFloat(arg.(float64), 'f', 2, 32)
	case string:
		return arg.(string)
	default:
		return ""
	}
}

func ToInt64(value interface{}) int64 {
	val := reflect.ValueOf(value)
	var d int64
	var err error
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	case float32, float64:
		d = int64(val.Float())
	case string:
		e, err1 := strconv.ParseFloat(val.String(), 64)
		if err1 != nil {
			d, err = strconv.ParseInt(val.String(), 10, 64)
		} else {
			d = int64(e)
		}
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	if err != nil {
		d = 0
	}
	return d
}

func ToInt(value interface{}) int {
	val := reflect.ValueOf(value)
	var d int
	var err error
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = int(val.Int())
	case uint, uint8, uint16, uint32, uint64:
		d = int(val.Uint())
	case float32, float64:
		d = int(val.Float())
	case string:
		e, err1 := strconv.ParseFloat(val.String(), 32)
		if err1 != nil {
			d, err = strconv.Atoi(val.String())
		} else {
			d = int(e)
		}

	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	if err != nil {
		log.Println(err)
		d = 0
	}

	return d
}

func MkKey(array []interface{}, key string, keyOnly bool) map[string]interface{} {
	data := map[string]interface{}{}

	if len(array) == 0 {
		return data
	}
	for _, v := range array {
		tmpMap := reflect.ValueOf(v)
		if tmpMap.Kind() != reflect.Map {
			return data
		}
		k := tmpMap.MapIndex(reflect.ValueOf(key))
		if !k.IsValid() {
			continue
		}
		kk := ToString(k.Interface())
		if keyOnly {
			data[kk] = kk
		} else {
			data[kk] = v
		}
	}
	return data
}

func MkKeyValue(array []interface{}, key string, value string) map[string]interface{} {
	data := map[string]interface{}{}

	if len(array) == 0 {
		return data
	}
	for _, v := range array {
		tmpMap := reflect.ValueOf(v)
		if tmpMap.Kind() != reflect.Map {
			return data
		}
		k := tmpMap.MapIndex(reflect.ValueOf(key))
		val := tmpMap.MapIndex(reflect.ValueOf(value))
		if !k.IsValid() {
			continue
		}
		if !val.IsValid() {
			continue
		}
		kk := ToString(k.Interface())
		vv := ToString(val.Interface())
		data[kk] = vv

	}
	return data
}

func MkMaxKeyValue(array []interface{}, key string, value string, maxKey string) map[string]interface{} {
	data := map[string]interface{}{}
	values := []string{}
	values = strings.Split(value, ",")
	if len(array) == 0 {
		return data
	}
	maxValArr := map[string]interface{}{}
	for _, v := range array {
		tmpMap := reflect.ValueOf(v)
		if tmpMap.Kind() != reflect.Map {
			return data
		}
		k := tmpMap.MapIndex(reflect.ValueOf(key))
		maxKeyVal := tmpMap.MapIndex(reflect.ValueOf(maxKey))

		if !k.IsValid() {
			continue
		}
		kk := ToString(k.Interface())
		mv := ToInt(maxKeyVal.Interface())
		if maxValArr[kk] == nil {
			maxValArr[kk] = mv
			vS := ""
			for _, vStr := range values {
				val := tmpMap.MapIndex(reflect.ValueOf(vStr))
				if !val.IsValid() {
					continue
				}
				vv := ToString(val.Interface())
				vS += vStr + ":" + vv + " "
			}
			data[kk] = vS
		}
		if ToInt(maxValArr[kk]) < mv {
			maxValArr[kk] = ToString(mv)
			vS := ""
			for _, vStr := range values {
				val := tmpMap.MapIndex(reflect.ValueOf(vStr))
				if !val.IsValid() {
					continue
				}
				vv := ToString(val.Interface())
				vS += "||*** " + vStr + ":" + vv + " ***|| "
			}
			data[kk] = vS
		}
	}
	return data
}

func GetKey(array map[string]interface{}) []string {
	data := []string{}

	if len(array) == 0 {
		return data
	}

	for key, _ := range array {
		data = append(data, key)
	}

	return data
}

/**
* array数据 格式化子元素的时间
* fields - 需要格式化为时间的字段名，多个字段用逗号隔开
* format - 时间格式
 */
func DataTimeFormat(array []interface{}, fields string, format string) []interface{} {
	fieldsMap := map[string]string{}
	fieldsArr := strings.Split(fields, ",")
	for _, f := range fieldsArr {
		fieldsMap[ToString(f)] = ToString(f)
	}
	if format == "" {
		format = "2006-01-02 15:04:05" // 时间不能修改，golang官方文档规定必须用这个时间点
	}
	data := []interface{}{}
	for _, val := range array {
		refVal := reflect.ValueOf(val)
		if refVal.Kind() != reflect.Map {
			continue
		}
		vv := map[string]interface{}{}
		for _, key := range refVal.MapKeys() {
			keyStr := key.String()
			valStr := refVal.MapIndex(key).String()
			if fieldsMap[keyStr] != "" {
				vv[keyStr] = time.Unix(ToInt64(valStr), 0).Format(format)
			} else {
				vv[keyStr] = valStr
			}
		}

		data = append(data, vv)
	}
	return data
}

// 处理http提交的二维数组
func ParseFormCollection(r *http.Request, typeName string, mustHave string) []map[string]string {
	var result []map[string]string
	r.ParseForm()
	for key, values := range r.Form {
		re := regexp.MustCompile(typeName + "\\[\\]\\[([a-zA-Z]+)\\]")
		matches := re.FindStringSubmatch(key)

		if len(matches) == 2 {
			for k, v := range values {
				for k >= len(result) {
					result = append(result, map[string]string{})
				}
				result[k][matches[1]] = strings.TrimSpace(v)
			}
		}
	}
	var data []map[string]string
	fieldsMap := map[string]string{}
	fieldsArr := strings.Split(mustHave, ",")
	for _, f := range fieldsArr {
		fieldsMap[ToString(f)] = ToString(f)
	}
	for _, values := range result {
		appendTo := true
		for _, mustKey := range fieldsMap {
			if values[mustKey] == "" {
				appendTo = false
				break
			}
		}
		if appendTo {
			data = append(data, values)
		}
	}
	return data
}

// 客户端IP
func GetIPAdress(r *http.Request) string {
	var ipAddress string
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			ip = strings.TrimSpace(ip)
			realIP := net.ParseIP(ip)
			if realIP != nil {
				ipAddress = ip
			}
		}
	}
	if len(ipAddress) == 0 {
		ipAddress, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ipAddress
}

func FileAppend(logFile string, data string) {
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(data)
	f.Sync()
}

// 字节转字符串
func Byte2String(bytes int64, unit string, returnUnit bool) string {
	units := [6]string{"B", "K", "M", "G", "T", "P"}
	offset := 0
	for bytes >= 1024 && offset < 5 {
		if unit != "" && units[offset] == strings.ToUpper(unit) {
			break
		}
		bytes = bytes >> 10
		offset++
	}
	result := ToString(bytes)
	if returnUnit == true {
		return result + units[offset]
	}
	return result
}

func ArrayDiff(array1 []string, array2 []string) []string {
	arr1 := map[string]string{}
	arr2 := map[string]string{}
	for _, val := range array1 {
		arr1[val] = val
	}
	for _, val := range array2 {
		arr2[val] = val
	}
	result := []string{}
	for k, v := range arr1 {
		if arr2[k] == "" {
			result = append(result, v)
		}
	}
	return result
}

// 适用于无法判断data为map类型的结构
func GetMapValueByKey(data interface{}, key string) string {
	refVal := reflect.ValueOf(data)
	if refVal.Kind() != reflect.Map {
		return ""
	}
	for _, k := range refVal.MapKeys() {
		keyStr := k.String()
		valStr := refVal.MapIndex(k).Interface()
		if keyStr == key {
			return ToString(valStr)
		}
	}
	return ""
}

// 检查文件名格式 小写英文、数字、连字符、下划线、点
func CheckFileName(file []string) (bool, string) {
	if len(file) > 0 {
		for _, f := range file {
			re := regexp.MustCompile("(^[a-z_0-9\\-\\.]{1,256}$)")
			matches := re.MatchString(f)
			if !matches {
				return false, f
			}
		}
		return true, ""
	}
	return false, ""
}

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func EncodeMD5(s string) string {
	dst, err := md5.Encrypt(s)
	if err != nil {
		return ""
	}
	return dst
}
