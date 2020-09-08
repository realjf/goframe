package conv

import (
	"encoding/json"
	"fmt"
	"github.com/realjf/goframe/encoding/binary"
	"reflect"
	"strconv"
	"strings"
)

var (
	// Empty strings.
	emptyStringMap = map[string]struct{}{
		"":      {},
		"0":     {},
		"no":    {},
		"off":   {},
		"false": {},
	}
)

type apiString interface {
	String() string
}

type apiError interface {
	Error() string
}


func Byte(v interface{}) byte {
	if val, ok := v.(byte); ok {
		return val
	}

	return Uint8(v)
}

func Bytes(v interface{}) []byte {
	if v == nil {
		return nil
	}
	switch value := v.(type) {
	case string:
		return []byte(value)
	case []byte:
		return value
	default:
		return binary.Encode(v)
	}
}

func Rune(v interface{}) rune {
	if val, ok := v.(rune); ok {
		return val
	}

	return rune(Int32(v))
}

func Runes(v interface{}) []rune {
	if val, ok := v.([]rune); ok {
		return val
	}

	return []rune(String(v))
}

func String(v interface{}) string {
	if v == nil {
		return ""
	}

	switch value := v.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case []rune:
		return string(value)
	default:
		if f, ok := value.(apiString); ok {
			return f.String()
		} else if f, ok := value.(apiError); ok {
			return f.Error()
		}else {
			if jsonContent, err := json.Marshal(value); err != nil {
				return fmt.Sprint(value)
			}else{
				return string(jsonContent)
			}
		}
	}
}

func Bool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch value := v.(type) {
	case bool:
		return value
	case []byte:
		if _, ok := emptyStringMap[strings.ToLower(string(value))]; ok {
			return false
		}
		return true
	case string:
		if _, ok := emptyStringMap[strings.ToLower(string(value))]; ok {
			return false
		}
		return true
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr:
			return !rv.IsNil()
		case reflect.Map:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			return rv.Len() != 0
		case reflect.Struct:
			return true
		default:
			s := strings.ToLower(String(v))
			if _, ok := emptyStringMap[s]; ok {
				return false
			}
			return true
		}
	}
}

func Int(v interface{}) int {
	if v == nil {
		return 0
	}
	if val, ok := v.(int); ok {
		return val
	}
	return int(Int64(v))
}

func Int8(v interface{}) int8 {
	if v == nil {
		return 0
	}
	if val, ok := v.(int8); ok {
		return val
	}
	return int8(Int64(v))
}

func Int16(v interface{}) int16 {
	if v == nil {
		return 0
	}
	if val, ok := v.(int16); ok {
		return val
	}
	return int16(Int64(v))
}

func Int32(v interface{}) int32 {
	if v == nil {
		return 0
	}
	if val, ok := v.(int32); ok {
		return val
	}
	return int32(Int64(v))
}

func Int64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch value := v.(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return binary.DecodeToInt64(value)
	default:
		s := String(value)
		// 十六进制
		if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
			if val, e := strconv.ParseInt(s[2:], 16, 64); e == nil {
				return val
			}
		}
		// 八进制
		if len(s) > 1 && s[0] == '0' {
			if val, e := strconv.ParseInt(s[1:], 8, 64); e == nil {
				return val
			}
		}
		// 十进制
		if val, e := strconv.ParseInt(s,10, 64); e == nil {
			return val
		}
		// 64位浮点数
		return int64(Float64(value))
	}
}

func Uint(v interface{}) uint {
	if v == nil {
		return 0
	}
	if val, ok := v.(uint); ok {
		return val
	}
	return uint(Uint64(v))
}

func Uint8(v interface{}) uint8 {
	if v == nil {
		return 0
	}
	if val, ok := v.(uint8); ok {
		return val
	}
	return uint8(Uint64(v))
}

func Uint16(v interface{}) uint16 {
	if v == nil {
		return 0
	}
	if val, ok := v.(uint16); ok {
		return val
	}
	return uint16(Uint64(v))
}

func Uint32(v interface{}) uint32 {
	if v == nil {
		return 0
	}
	if val, ok := v.(uint32); ok {
		return val
	}
	return uint32(Uint64(v))
}

func Uint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}
	switch value := v.(type) {
	case int:
		return uint64(value)
	case int8:
		return uint64(value)
	case int16:
		return uint64(value)
	case int32:
		return uint64(value)
	case int64:
		return uint64(value)
	case uint:
		return uint64(value)
	case uint8:
		return uint64(value)
	case uint16:
		return uint64(value)
	case uint32:
		return uint64(value)
	case uint64:
		return value
	case float32:
		return uint64(value)
	case float64:
		return uint64(value)
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return binary.DecodeToUint64(value)
	default:
		s := String(value)
		// 十六进制
		if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
			if val, e := strconv.ParseUint(s[2:], 16, 64); e == nil {
				return val
			}
		}
		// 八进制
		if len(s) > 1 && s[0] == '0' {
			if val, e := strconv.ParseUint(s[1:], 8, 64); e == nil {
				return val
			}
		}
		// 十进制
		if val, e := strconv.ParseUint(s, 10, 64); e == nil {
			return val
		}
		// float64
		return uint64(Float64(value))
	}
}

func Float32(v interface{}) float32 {
	if v == nil {
		return 0
	}
	switch value := v.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case []byte:
		return binary.DecodeToFloat32(value)
	default:
		val, _ := strconv.ParseFloat(String(v), 64)
		return Float32(val)
	}
}

func Float64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch value := v.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case []byte:
		return binary.DecodeToFloat64(value)
	default:
		val, _ := strconv.ParseFloat(String(v), 64)
		return val
	}
}