package validate

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type SupportedEnvTypes interface {
	~int | ~float64 | ~bool | ~string
}

func coerceString[T SupportedEnvTypes](val string) (T, error) {
	var result T

	switch any(result).(type) {
	case int:
		var tmp int
		tmp, err := strconv.Atoi(val)
		if err != nil {
			return result, err
		}

		result = any(tmp).(T)
	case float64:
		tmp, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return result, err
		}
		result = any(tmp).(T)

	case bool:
		tmp, err := strconv.ParseBool(val)
		if err != nil {
			return result, err
		}
		result = any(tmp).(T)
	case string:
		result = any(val).(T)
	default:
		return result, fmt.Errorf("unsupported type: %T", result)
	}

	return result, nil
}

func isZeroValue(x any) bool {
	if x == nil {
		return true
	}

	v := reflect.ValueOf(x)
	if !v.IsValid() {
		return true
	}

	// Check if the value is the zero value for its type
	zeroValue := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zeroValue.Interface())
}

func Env[T SupportedEnvTypes](key string, rulesSets []RuleSet, defaultValue ...T) T {

	str := os.Getenv(key)

	val, err := coerceString[T](str)

	if err != nil || isZeroValue(val) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		} else {
			panic(fmt.Errorf("failed to parse env %s: %v", key, err))
		}
	}

	fieldName := key
	fieldValue := val

	for _, set := range rulesSets {
		set.FieldValue = fieldValue
		set.FieldName = fieldName
		if !set.ValidateFunc(set) {
			msg := set.MessageFunc(set)
			panic(fmt.Sprintf("Error parsing env %s: %s", key, msg))
		}
	}

	return val
}
