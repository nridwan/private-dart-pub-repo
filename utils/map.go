package utils

import "reflect"

func StructToMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name
		jsonTag := val.Type().Field(i).Tag.Get("json")
		if jsonTag != "" {
			fieldName = jsonTag
		}
		result[fieldName] = field.Interface()
	}
	return result
}
