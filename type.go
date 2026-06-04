package goutil

import (
	"reflect"
	"strconv"
)

type Hashable interface {
	string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128
}

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64
}

type NullType[T any] struct {
	Null T
}

type ToInterface struct {
	Val interface{}
}

// IsZeroOfUnderlyingType can be used to determine if an interface{} in null or empty
func IsZeroOfUnderlyingType(x any) bool {
	// return x == nil || x == reflect.Zero(reflect.TypeOf(x)).Interface()
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// toString converts multiple types to a string|[]byte
//
// accepts: string, []byte, byte, int (and variants), [][]byte, []interface{}
func toString[T interface{ string | []byte }](val any) T {
	switch val.(type) {
	case string:
		return T(val.(string))
	case []byte:
		return T(val.([]byte))
	case byte:
		return T([]byte{val.(byte)})
	case int:
		return T(strconv.Itoa(val.(int)))
	case int64:
		return T(strconv.Itoa(int(val.(int64))))
	case int32:
		return T(strconv.Itoa(int(val.(int32))))
	case int16:
		return T([]byte{byte(val.(int16))})
	case int8:
		return T([]byte{byte(val.(int8))})
	case uintptr:
		return T(strconv.FormatUint(uint64(val.(uintptr)), 10))
	case uint:
		return T(strconv.FormatUint(uint64(val.(uint)), 10))
	case uint64:
		return T(strconv.FormatUint(val.(uint64), 10))
	case uint32:
		return T(strconv.FormatUint(uint64(val.(uint32)), 10))
	case uint16:
		return T(strconv.FormatUint(uint64(val.(uint16)), 10))
	case float64:
		return T(strconv.FormatFloat(val.(float64), 'f', -1, 64))
	case float32:
		return T(strconv.FormatFloat(float64(val.(float32)), 'f', -1, 32))
	case []interface{}:
		b := make([]byte, len(val.([]interface{})))
		for i, v := range val.([]interface{}) {
			b[i] = byte(toNumber[int32](v))
		}
		return T(b)
	case []int:
		b := make([]byte, len(val.([]int)))
		for i, v := range val.([]int) {
			b[i] = byte(v)
		}
		return T(b)
	case []int64:
		b := make([]byte, len(val.([]int64)))
		for i, v := range val.([]int64) {
			b[i] = byte(v)
		}
		return T(b)
	case []int32:
		b := make([]byte, len(val.([]int32)))
		for i, v := range val.([]int32) {
			b[i] = byte(v)
		}
		return T(b)
	case []int16:
		b := make([]byte, len(val.([]int16)))
		for i, v := range val.([]int16) {
			b[i] = byte(v)
		}
		return T(b)
	case []int8:
		b := make([]byte, len(val.([]int8)))
		for i, v := range val.([]int8) {
			b[i] = byte(v)
		}
		return T(b)
	case []uint:
		b := make([]byte, len(val.([]uint)))
		for i, v := range val.([]uint) {
			b[i] = byte(v)
		}
		return T(b)
	case []uint16:
		b := make([]byte, len(val.([]uint16)))
		for i, v := range val.([]uint16) {
			b[i] = byte(v)
		}
		return T(b)
	case []uint32:
		b := make([]byte, len(val.([]uint32)))
		for i, v := range val.([]uint32) {
			b[i] = byte(v)
		}
		return T(b)
	case []uint64:
		b := make([]byte, len(val.([]uint64)))
		for i, v := range val.([]uint64) {
			b[i] = byte(v)
		}
		return T(b)
	case []uintptr:
		b := make([]byte, len(val.([]uintptr)))
		for i, v := range val.([]uintptr) {
			b[i] = byte(v)
		}
		return T(b)
	case []string:
		b := []byte{}
		for _, v := range val.([]string) {
			b = append(b, []byte(v)...)
		}
		return T(b)
	case [][]byte:
		b := []byte{}
		for _, v := range val.([][]byte) {
			b = append(b, v...)
		}
		return T(b)
	default:
		return T("")
	}
}

// toNumber converts multiple types to a number
//
// accepts: int (and variants), string, []byte, byte, bool
func toNumber[T interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float64 | float32
}](val any) T {
	switch val.(type) {
	case int:
		return T(val.(int))
	case int64:
		return T(val.(int64))
	case float64:
		return T(val.(float64))
	case float32:
		return T(val.(float32))
	case string:
		var varT interface{} = T(0)
		if _, ok := varT.(float64); ok {
			if f, err := strconv.ParseFloat(val.(string), 64); err == nil {
				return T(f)
			}
		} else if _, ok := varT.(float32); ok {
			if f, err := strconv.ParseFloat(val.(string), 32); err == nil {
				return T(f)
			}
		} else if i, err := strconv.Atoi(val.(string)); err == nil {
			return T(i)
		}
		return 0
	case []byte:
		if i, err := strconv.Atoi(string(val.([]byte))); err == nil {
			return T(i)
		}
		return 0
	case byte:
		if i, err := strconv.Atoi(string(val.(byte))); err == nil {
			return T(i)
		}
		return 0
	case bool:
		if val.(bool) {
			return 1
		}
		return 0
	case int8:
		return T(val.(int8))
	case int16:
		return T(val.(int16))
	case int32:
		return T(val.(int32))
	case uint:
		return T(val.(uint))
	case uint16:
		return T(val.(uint16))
	case uint32:
		return T(val.(uint32))
	case uint64:
		return T(val.(uint64))
	case uintptr:
		return T(val.(uintptr))
	default:
		return 0
	}
}

// SupportedType is an interface containing the types which are supported by the ToType method
type SupportedType interface {
	string | []byte | byte | bool |
		int | int64 | int32 | int16 | int8 |
		uint | uint64 | uint32 | uint16 | /* uint8 | */ uintptr |
		float64 | float32
}

// ToType attempts to converts an interface{} from the many possible types in golang, to a specific type of your choice
//
// if it fails to convert, it will return a nil/zero value for the appropriate type
func ToType[T SupportedType](val any) T {
	var varT interface{} = NullType[T]{}.Null
	switch varT.(type) {
	// basic
	case string:
		return ToInterface{toString[string](val)}.Val.(T)
	case []byte:
		return ToInterface{toString[[]byte](val)}.Val.(T)
	case byte:
		if b := toString[[]byte](val); len(b) != 0 {
			return ToInterface{b[0]}.Val.(T)
		}
		return ToInterface{byte(0)}.Val.(T)
	case bool:
		return ToInterface{!IsZeroOfUnderlyingType(val)}.Val.(T)

	// int
	case int:
		return ToInterface{toNumber[int](val)}.Val.(T)
	case int64:
		return ToInterface{toNumber[int64](val)}.Val.(T)
	case int32:
		return ToInterface{toNumber[int32](val)}.Val.(T)
	case int16:
		return ToInterface{toNumber[int16](val)}.Val.(T)
	case int8:
		return ToInterface{toNumber[int8](val)}.Val.(T)

	// uint
	case uintptr:
		return ToInterface{toNumber[uintptr](val)}.Val.(T)
	case uint:
		return ToInterface{toNumber[uint](val)}.Val.(T)
	case uint64:
		return ToInterface{toNumber[uint64](val)}.Val.(T)
	case uint32:
		return ToInterface{toNumber[uint32](val)}.Val.(T)
	case uint16:
		return ToInterface{toNumber[uint16](val)}.Val.(T)

	// float
	case float64:
		return ToInterface{toNumber[float64](val)}.Val.(T)
	case float32:
		return ToInterface{toNumber[float32](val)}.Val.(T)

	default:
		return NullType[T]{}.Null
	}
}

/* func ToType[T SupportedType](val interface{}) T {
	// basic
	var varT interface{} = ""
	if _, ok := varT.(T); ok {
		return ToInterface{toString[string](val)}.Val.(T)
	}

	varT = []byte{}
	if _, ok := varT.(T); ok {
		return ToInterface{toString[[]byte](val)}.Val.(T)
	}

	varT = byte(0)
	if _, ok := varT.(T); ok {
		if b := toString[[]byte](val); len(b) != 0 {
			return ToInterface{b[0]}.Val.(T)
		}
		return ToInterface{byte(0)}.Val.(T)
	}

	varT = rune(0)
	if _, ok := varT.(T); ok {
		if b := toString[[]byte](val); len(b) != 0 {
			return ToInterface{rune(b[0])}.Val.(T)
		}
		return ToInterface{rune(0)}.Val.(T)
	}

	varT = false
	if _, ok := varT.(T); ok {
		return ToInterface{!IsZeroOfUnderlyingType(val)}.Val.(T)
	}

	// int
	varT = int(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int](val)}.Val.(T)
	}

	varT = int64(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int64](val)}.Val.(T)
	}

	varT = int32(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int32](val)}.Val.(T)
	}

	varT = int16(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int16](val)}.Val.(T)
	}

	varT = int8(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int8](val)}.Val.(T)
	}

	// uint
	varT = uintptr(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uintptr](val)}.Val.(T)
	}

	varT = uint(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint](val)}.Val.(T)
	}

	varT = uint64(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint64](val)}.Val.(T)
	}

	varT = uint32(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint32](val)}.Val.(T)
	}

	varT = uint16(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint16](val)}.Val.(T)
	}

	varT = uint8(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint8](val)}.Val.(T)
	}

	// float
	varT = float64(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[float64](val)}.Val.(T)
	}

	varT = float32(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[float32](val)}.Val.(T)
	}

	return NullType[T]{}.Null
} */
