package structs

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ToMapStr convert struct to map[string][string] for redis hash, if tag is null, key is snake format
// tag: `map:"key,omitempty"`
func ToMapStr(in interface{}, fields ...string) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMapStr only accepts struct or struct pointer; got %T", v)
	}

	set := make(map[string]struct{}, len(fields))
	for _, s := range fields {
		set[s] = struct{}{}
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		mt := f.Tag.Get("map")
		var name string
		var omitempty bool
		if len(mt) == 0 {
			name = camelToSnake(f.Name)

		} else {
			split := strings.Split(mt, ",")
			name, omitempty = split[0], split[1] == "omitempty"
		}
		if len(set) == 0 {
			if omitempty && v.Field(i).IsZero() {
				continue
			}
			out[name] = fmt.Sprintf("%v", v.Field(i).Interface())
		} else {
			if _, ok := set[name]; ok {
				out[name] = fmt.Sprintf("%v", v.Field(i).Interface())
			}
		}
	}

	return out, nil
}

func FillStruct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		err := setField(obj, snakeToPascal(k), v)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	var err error
	if structFieldType != val.Type() {
		val, err = typeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

func typeConversion(value string, t string) (reflect.Value, error) {
	if t == "string" {
		return reflect.ValueOf(value), nil
	} else if t == "time.Time" || t == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if t == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if t == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if t == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if t == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if t == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if t == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	return reflect.ValueOf(value), fmt.Errorf("unknown type: %s", t)
}

func camelToSnake(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func snakeToPascal(s string) string {
	sb := new(strings.Builder)
	j := false
	sb.WriteString(strings.ToUpper(string(s[0])))
	for i := 1; i < len(s); i++ {
		d := string(s[i])
		if j == true {
			sb.WriteString(strings.ToUpper(d))
			j = false
			continue
		}
		if d == "_" {
			j = true
			continue
		}
		sb.WriteString(strings.ToLower(d))
	}
	return sb.String()
}
