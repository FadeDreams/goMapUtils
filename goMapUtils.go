package goMapUtils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func MapType(input map[interface{}]interface{}) (keyType, valueType string) {
	keyType = reflect.TypeOf(input).Key().String()
	valueType = reflect.TypeOf(input).Elem().String()
	return
}

func ConcatMaps(map1, map2 map[interface{}]interface{}) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})
	for key, value := range map1 {
		result[key] = value
	}
	for key, value := range map2 {
		result[key] = value
	}
	return result
}

func ContainsKey(m map[interface{}]interface{}, key interface{}) bool {
	_, ok := m[key]
	return ok
}

func Clone(o interface{}) interface{} {
	if o == nil {
		return nil
	}

	switch reflect.TypeOf(o).Kind() {
	case reflect.Map:
		return deepCopyMap(o.(map[interface{}]interface{}))
	case reflect.Slice:
		return deepCopySlice(o.([]interface{}))
	case reflect.Struct:
		return deepCopyStruct(o)
	default:
		return o
	}
}

func deepCopyMap(m map[interface{}]interface{}) map[interface{}]interface{} {
	copy := make(map[interface{}]interface{})
	for k, v := range m {
		copy[k] = Clone(v)
	}
	return copy
}

func deepCopySlice(s []interface{}) []interface{} {
	copy := make([]interface{}, len(s))
	for i, v := range s {
		copy[i] = Clone(v)
	}
	return copy
}

func deepCopyStruct(o interface{}) interface{} {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)

	copy := reflect.New(t).Elem()
	for i := 0; i < t.NumField(); i++ {
		copy.Field(i).Set(deepCopyValue(v.Field(i)))
	}
	return copy.Interface()
}

func deepCopyValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Map:
		return reflect.ValueOf(deepCopyMap(v.Interface().(map[interface{}]interface{})))
	case reflect.Slice:
		return reflect.ValueOf(deepCopySlice(v.Interface().([]interface{})))
	case reflect.Struct:
		return reflect.ValueOf(deepCopyStruct(v.Interface()))
	default:
		return v
	}
}

func ConvertJSONToMap(jsonString string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ConvertMapToJSON(data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func IsMapEmpty(data map[string]interface{}) bool {
	return len(data) == 0
}

func GetValueFromMap(o interface{}, t string) interface{} {
	if o == nil || t == "" || reflect.TypeOf(o).Kind() != reflect.Map {
		return nil
	}

	keyArr := strings.Split(t, ".")
	firstKey := keyArr[0]
	keyArr = keyArr[1:]

	if len(keyArr) == 0 {
		return o.(map[string]interface{})[firstKey]
	}

	return GetValueFromMap(o.(map[string]interface{})[firstKey], strings.Join(keyArr, "."))
}

func IsTheSameMap(o1, o2 interface{}, skipFields []interface{}) bool {
	if (o1 == nil && o2 != nil) || (o1 != nil && o2 == nil) {
		return false
	}

	if reflect.TypeOf(o1) != reflect.TypeOf(o2) {
		return false
	}

	if reflect.TypeOf(o1).Kind() == reflect.Slice && reflect.TypeOf(o1).Kind() != reflect.TypeOf(o2).Kind() {
		return false
	}

	if reflect.TypeOf(o1).Kind() == reflect.String {
		return strings.ToLower(strings.TrimSpace(o1.(string))) == strings.ToLower(strings.TrimSpace(o2.(string)))
	}

	if reflect.TypeOf(o1).Kind() == reflect.Slice {
		o1Slice, o2Slice := o1.([]interface{}), o2.([]interface{})
		if len(o1Slice) != len(o2Slice) {
			return false
		}
		for _, item := range o1Slice {
			if !contains(o2Slice, item.(string)) {
				return false
			}
		}
		return true
	}

	if reflect.TypeOf(o1).Kind() == reflect.Map {
		o1Map, o2Map := o1.(map[string]interface{}), o2.(map[string]interface{})
		if len(o1Map) != len(o2Map) {
			return false
		}
		for o1Key, o1Value := range o1Map {
			if !contains(skipFields, o1Key) {
				o2Value, ok := o2Map[o1Key]
				if !ok || !IsTheSameMap(o1Value, o2Value, skipFields) {
					return false
				}
			}
		}
		return true
	}

	return o1 == o2
}

func contains(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func SortMapByKeys(m map[interface{}]interface{}) []interface{} {
	keys := make([]interface{}, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		// Convert keys to strings and compare their lexicographic order
		return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
	})
	return keys
}

type KeyExtractor func(entry map[interface{}]interface{}) string

func SortMapByCustomKey(m map[interface{}]interface{}, keyExtractor KeyExtractor) []string {
	keys := make([]string, 0, len(m))
	for k, v := range m {
		keys = append(keys, keyExtractor(map[interface{}]interface{}{k: v}))
	}
	sort.Strings(keys)
	return keys
}

func PrettyPrintMap(m map[interface{}]interface{}) {
	// Iterate over the map and print each key-value pair
	for k, v := range m {
		fmt.Printf("%v: %v\n", k, v)
	}
}
