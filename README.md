Set of tools to handle Map type in Golang.

### Installation
```
  go get -u github.com/fadedreams/goMapUtils
```

#### note:
-  Examine the function definition, and if your map type differs from the [interface{}]interface{}, convert it first.

```go
// convert to type [interface{}]interface{}
map1 := map[string]int{"a": 1, "b": 2}
map2 := map[int]string{1: "c", 2: "d"}

map1Interface := make(map[interface{}]interface{})
for key, value := range map1 {
	map1Interface[key] = value
}
```

| Summary of methods |  |
| ------------- | ------------- |
| PrettyPrintMap |  ConcatMaps | 
| ContainsKey  | IsTheSameMap | 
| MapType | GetValueFromMap | 
| Clone | CloneAsync |
| ConvertMapToJSON |ConvertJSONToMap | 
| IsMapEmpty | SortMapByKeys | 
| SortMapByCustomKey | IterateMap |

### Usage examples

#### PrettyPrintMap
```go
m := make(map[interface{}]interface{})
m["foo"] = 1
m[123] = "bar"
m[true] = 3.14

PrettyPrintMap(m)
```

#### ContainsKey
```go
m := make(map[interface{}]interface{})
m["foo"] = 1
m[123] = "bar"

fmt.Println(ContainsKey(m, "foo"))  // prints "true"
fmt.Println(ContainsKey(m, 123))   // prints "true"
fmt.Println(ContainsKey(m, "baz")) // prints "false"
```

#### MapType
```go
input := map[string]int{"a": 1, "b": 2}
keyType, valueType := MapType(input)
fmt.Printf("Map has key type %q and value type %q\n", keyType, valueType)

```

#### Clone
    Clone is a utility function for creating a deep copy of an object.
```go
type Person struct {
    Name string
    Age  int
}

p1 := Person{
    Name: "John",
    Age:  30,
}

p2 := Clone(p1).(Person)

fmt.Println(p1) // {John 30}
fmt.Println(p2) // {John 30}

p2.Name = "Jane"

fmt.Println(p1) // {John 30}
fmt.Println(p2) // {Jane 30}
```

#### CloneAsync
The Clone function is an asynchronous deep copy function that creates a new copy of an input value by recursively copying all its child values asynchronously, using goroutines and the sync.WaitGroup type. It can handle input values of type map, slice, and struct, and returns the input value for any other type.

```go
import (
    "fmt"
    "reflect"
    "sync"
)

type User struct {
    ID   int
    Name string
}

func main() {
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

    copy := CloneAsync(data).(map[interface{}]interface{})

    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Copy: %v\n", copy)
    fmt.Println("Are original and copy equal? ", reflect.DeepEqual(data, copy))
}

```

#### ConcatMaps
```go
//if your maps type is not same as [interface{}]interface{} 
//convert them using before calling ConcatMaps
//for example
//map1 := map[string]int{"a": 1, "b": 2}
//map2 := map[int]string{1: "c", 2: "d"}

//map1Interface := make(map[interface{}]interface{})
//for key, value := range map1 {
//	map1Interface[key] = value
//}

//// Convert map2 to a map with interface{} key and value types
//map2Interface := make(map[interface{}]interface{})
//for key, value := range map2 {
//	map2Interface[key] = value
//}

map1 := map[interface{}]interface{}{"a": 1, "b": "hello"}
map2 := map[interface{}]interface{}{1: 3.14, 2: []int{1, 2, 3}}
concatenatedMap := ConcatMaps(map1, map2)
fmt.Println(concatenatedMap)  // Output: map[a:1 b:hello 1:3.14 2:[1 2 3]]

map1 = map[interface{}]interface{}{"a": 1, "b": 2}
map2 = map[interface{}]interface{}{"b": 3, "c": 4}
concatenatedMap = ConcatMaps(map1, map2)
fmt.Println(concatenatedMap)  // Output: map[a:1 b:3 c:4]

```


#### IsTheSameMap
```go
o1 := map[string]interface{}{
    "name": "John",
    "age": 30,
}

o2 := map[string]interface{}{
    "name": "john",
    "age": 30,
}

if IsTheSameMap(o1, o2, []interface{}{"name"}) {
    fmt.Println("The maps are the same.")
} else {
    fmt.Println("The maps are not the same.")
}
```

#### GetValueFromMap
```go
data = map[string]interface{}{}
value := GetValueFromMap(data, "a.b.c")
fmt.Println(value)  // Output: nil

data = map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 3.14, "d": []int{1, 2, 3}}}
value = GetValueFromMap(data, "a")
fmt.Println(value)  // Output: 1

value = GetValueFromMap(data, "b.c")
fmt.Println(value)  // Output: 3.14
```

#### ConvertJSONToMap & ConvertMapToJSON
```go
// Convert JSON string to a map
jsonString := `{"a": 1, "b": "hello"}`
data, err := ConvertJSONToMap(jsonString)
if err != nil {
  fmt.Println(err)
  return
}
fmt.Println(data)  // Output: map[a:1 b:hello]

// Convert map to a JSON string
jsonData, err := ConvertMapToJSON(data)
if err != nil {
  fmt.Println(err)
  return
}
fmt.Println(jsonData)  // Output: {"a":1,"b":"hello"}
```


#### IsMapEmpty
```go
// Check if map is empty
data := map[string]interface{}{"a": 1, "b": "hello"}
fmt.Println(IsMapEmpty(data))  // Output: false

data = map[string]interface{}{}
fmt.Println(IsMapEmpty(data))  // Output: true
```

#### SortMapByKeys

```go
myMap := map[interface{}]interface{}{
	"foo": 1,
	"bar": 2,
	"baz": 3,
}
sortedKeys := sortMapByKeys(myMap)
fmt.Println(sortedKeys) // Output: [bar baz foo]
```

#### SortMapByCustomKey
```go
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
sortedKeys := sortMapByCustomKey(myMap, keyExtractor)
fmt.Println(sortedKeys) // prints ["four", "one", "three", "two"]
```

#### IterateMap
This function takes a map as an argument and returns a channel that allows the caller to iterate over the map by receiving pairs of key and value from the channel.
```go
m := map[interface{}]interface{}{
  "a": 1,
  "b": 2,
  "c": 3,
  4:   "four",
}

for pair := range IterateMap(m) {
  key := pair[0]
  value := pair[1]
  fmt.Println("Key:", key, "Value:", value)
}
```
