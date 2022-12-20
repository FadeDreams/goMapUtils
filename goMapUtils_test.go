package goMapUtils

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestConcatMaps(t *testing.T) {
	// Test concatenating two empty maps
	map1 := map[interface{}]interface{}{}
	map2 := map[interface{}]interface{}{}
	expectedResult := map[interface{}]interface{}{}
	result := ConcatMaps(map1, map2)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("ConcatMaps(%v, %v) = %v, expected %v", map1, map2, result, expectedResult)
	}

	// Test concatenating two maps with different types
	map1 = map[interface{}]interface{}{"a": 1, "b": "hello"}
	map2 = map[interface{}]interface{}{1: 3.14, 2: []int{1, 2, 3}}
	expectedResult = map[interface{}]interface{}{"a": 1, "b": "hello", 1: 3.14, 2: []int{1, 2, 3}}
	result = ConcatMaps(map1, map2)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("ConcatMaps(%v, %v) = %v, expected %v", map1, map2, result, expectedResult)
	}

	// Test concatenating two maps with overlapping keys
	map1 = map[interface{}]interface{}{"a": 1, "b": 2}
	map2 = map[interface{}]interface{}{"b": 3, "c": 4}
	expectedResult = map[interface{}]interface{}{"a": 1, "b": 3, "c": 4}
	result = ConcatMaps(map1, map2)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("ConcatMaps(%v, %v) = %v, expected %v", map1, map2, result, expectedResult)
	}
}

func TestContainsKey(t *testing.T) {
	m := make(map[interface{}]interface{})
	m["foo"] = 1
	m[123] = "bar"

	// Test for keys that are present in the map
	if !ContainsKey(m, "foo") {
		t.Error("Expected ContainsKey to return true for key 'foo', got false")
	}
	if !ContainsKey(m, 123) {
		t.Error("Expected ContainsKey to return true for key 123, got false")
	}

	// Test for keys that are not present in the map
	if ContainsKey(m, "baz") {
		t.Error("Expected ContainsKey to return false for key 'baz', got true")
	}
	if ContainsKey(m, 456) {
		t.Error("Expected ContainsKey to return false for key 456, got true")
	}
}

func compareMaps(x, y map[string]interface{}) bool {
	if len(x) != len(y) {
		return false
	}

	for k, v := range x {
		if !compareValues(v, y[k]) {
			return false
		}
	}

	return true
}

func compareValues(x, y interface{}) bool {
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)

	if v1.Kind() != v2.Kind() {
		return false
	}

	switch v1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v1.Int() == v2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v1.Uint() == v2.Uint()
	case reflect.Float32, reflect.Float64:
		return v1.Float() == v2.Float()
	case reflect.String:
		return v1.String() == v2.String()
	case reflect.Bool:
		return v1.Bool() == v2.Bool()
	case reflect.Slice:
		if v1.Len() != v2.Len() {
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if !compareValues(v1.Index(i).Interface(), v2.Index(i).Interface()) {
				return false
			}
		}
		return true
	case reflect.Map:
		if v1.Len() != v2.Len() {
			return false
		}
		for _, k := range v1.MapKeys() {
			if !compareValues(v1.MapIndex(k).Interface(), v2.MapIndex(k).Interface()) {
				return false
			}
		}
		return true
	default:
		return v1.Interface() == v2.Interface()
	}
}

func TestConvertJSONToMap(t *testing.T) {
	jsonString := `{"name": "John", "age": 30, "skills": ["programming", "writing"]}`
	expectedData := map[string]interface{}{
		"name": "John",
		"age":  30,
		"skills": []interface{}{
			"programming",
			"writing",
		},
	}

	data, err := ConvertJSONToMap(jsonString)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if compareMaps(data, expectedData) {
		t.Errorf("Expected %v, got %v", expectedData, data)
	}

}

func TestConvertMapToJSON(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	expectedJSON := `{"key1":"value1","key2":123,"key3":true}`
	jsonData, err := ConvertMapToJSON(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if jsonData != expectedJSON {
		t.Errorf("Expected %q but got %q", expectedJSON, jsonData)
	}
}

func TestIsMapEmpty(t *testing.T) {
	// Test an empty map
	data := make(map[string]interface{})
	if !IsMapEmpty(data) {
		t.Errorf("Expected isMapEmpty to return true for an empty map but got false")
	}

	// Test a non-empty map
	data["key"] = "value"
	if IsMapEmpty(data) {
		t.Errorf("Expected isMapEmpty to return false for a non-empty map but got true")
	}
}

func TestGetValueFromMap(t *testing.T) {
	// Test with a simple object
	obj := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	val := GetValueFromMap(obj, "key1")
	if val != "value1" {
		t.Errorf("Expected value %q but got %q", "value1", val)
	}

	val = GetValueFromMap(obj, "key2")
	if val != 123 {
		t.Errorf("Expected value %d but got %v", 123, val)
	}

	val = GetValueFromMap(obj, "key3")
	if val != true {
		t.Errorf("Expected value %t but got %v", true, val)
	}

	// Test with a nested object
	obj = map[string]interface{}{
		"key1": map[string]interface{}{
			"key2": "value2",
		},
	}

	val = GetValueFromMap(obj, "key1.key2")
	if val != "value2" {
		t.Errorf("Expected value %q but got %q", "value2", val)
	}

	// Test with an invalid key
	val = GetValueFromMap(obj, "invalid")
	if val != nil {
		t.Errorf("Expected value %v but got %v", nil, val)
	}
}

func TestisTheSameMap(t *testing.T) {
	// Generate random values for the objects
	rand.Seed(time.Now().UnixNano())
	numKeys := rand.Intn(5) + 1
	obj1 := make(map[string]interface{})
	obj2 := make(map[string]interface{})
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key%d", i+1)
		obj1[key] = rand.Intn(1000)
		obj2[key] = obj1[key]
	}

	// Check that the objects are considered the same
	if !IsTheSameMap(obj1, obj2, nil) {
		t.Errorf("Expected isTheSameMap to return true but got false")
	}

	// Modify one of the objects and check that they are not considered the same
	obj2["key1"] = obj1["key1"].(int) + 1
	if IsTheSameMap(obj1, obj2, nil) {
		t.Errorf("Expected isTheSameMap to return false but got true")
	}
}

// clone test
func TestClone(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{
			input:    map[interface{}]interface{}{1: "one", 2: "two"},
			expected: map[interface{}]interface{}{1: "one", 2: "two"},
		},
		{
			input:    []interface{}{1, 2, 3},
			expected: []interface{}{1, 2, 3},
		},
		{
			input: TestStruct{
				Name: "Alice",
				Age:  20,
			},
			expected: TestStruct{
				Name: "Alice",
				Age:  20,
			},
		},
		{
			input:    "hello",
			expected: "hello",
		},
		{
			input:    nil,
			expected: nil,
		},
	}

	for _, test := range tests {
		result := Clone(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test failed: expected %v, got %v", test.expected, result)
		}
	}
}

func TestSortMapByKeys(t *testing.T) {
	testCases := []struct {
		input    map[interface{}]interface{}
		expected []interface{}
	}{
		{
			input: map[interface{}]interface{}{
				"foo": 1,
				"bar": 2,
				"baz": 3,
			},
			expected: []interface{}{"bar", "baz", "foo"},
		},
		{
			input: map[interface{}]interface{}{
				"baz": 3,
				"foo": 1,
				"bar": 2,
			},
			expected: []interface{}{"bar", "baz", "foo"},
		},
		{
			input: map[interface{}]interface{}{
				1: "foo",
				2: "bar",
				3: "baz",
			},
			expected: []interface{}{1, 2, 3},
		},
		{
			input: map[interface{}]interface{}{
				3: "baz",
				1: "foo",
				2: "bar",
			},
			expected: []interface{}{1, 2, 3},
		},
		{
			input:    map[interface{}]interface{}{},
			expected: []interface{}{},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.input), func(t *testing.T) {
			output := SortMapByKeys(tc.input)
			if !reflect.DeepEqual(output, tc.expected) {
				t.Errorf("Got %v, expected %v", output, tc.expected)
			}
		})
	}
}

func TestSortMapByCustomKey(t *testing.T) {
	// create a map of integers to strings
	myMap := map[interface{}]interface{}{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
	}

	// define a KeyExtractor function that extracts the value of the entry as the key
	keyExtractor := func(entry map[interface{}]interface{}) string {
		for _, v := range entry {
			return v.(string)
		}
		return ""
	}

	// sort the map by the values of the entries
	sortedKeys := SortMapByCustomKey(myMap, keyExtractor)

	// check that the sorted keys are as expected
	expectedKeys := []string{"four", "one", "three", "two"}
	if !reflect.DeepEqual(sortedKeys, expectedKeys) {
		t.Errorf("Expected sorted keys to be %v, got %v", expectedKeys, sortedKeys)
	}
}

func TestPrettyPrintMap(t *testing.T) {
	// Capture the output of the prettyPrintMap function
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a test map and pass it to the prettyPrintMap function
	m := make(map[interface{}]interface{})
	m["foo"] = 1
	m[123] = "bar"
	m[true] = 3.14
	PrettyPrintMap(m)

	// Restore the standard output and read the captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Check if the output is what we expected
	expectedOutput := "foo: 1\n123: bar\ntrue: 3.14\n"
	if buf.String() != expectedOutput {
		t.Errorf("Expected prettyPrintMap to output '%s', got '%s'", expectedOutput, buf.String())
	}
}

func TestIterateMap(t *testing.T) {
	m := map[interface{}]interface{}{
		"a": 1,
		2:   "b",
	}

	expected := [][2]interface{}{
		{"a", 1},
		{2, "b"},
	}

	i := 0
	for pair := range IterateMap(m) {
		if pair != expected[i] {
			t.Errorf("Expected %v, got %v", expected[i], pair)
		}
		i++
	}
}

func TestCloneAsync(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	// Create a sample data structure to clone
	data := map[interface{}]interface{}{
		"users": []interface{}{
			User{ID: 1, Name: "Alice"},
			User{ID: 2, Name: "Bob"},
		},
		"settings": map[interface{}]interface{}{
			"color": "blue",
			"theme": "light",
		},
	}

	// Clone the data structure asynchronously
	copy := CloneAsync(data).(map[interface{}]interface{})

	// Check that the original and copy are equal
	if !reflect.DeepEqual(data, copy) {
		t.Errorf("Expected copy to be equal to original, but got %v", copy)
	}
}
