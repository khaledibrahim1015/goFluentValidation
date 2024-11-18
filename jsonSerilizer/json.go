// Building a Simple Serialization/Deserialization Mechanism
package jsonserilizer

import (
	"fmt"
	"reflect"
)

// Result is a map type used for serialization/deserialization.
type Result map[string]interface{}

// Serialize converts a struct into a Result map.
func Serialize(s interface{}) Result {
	result := make(Result)
	// assume s will be a struct
	rTyp := reflect.TypeOf(s)
	rVal := reflect.ValueOf(s)

	// Ensure the input is a struct.
	if rTyp.Kind() != reflect.Struct {
		return nil
	}

	// Iterate over the struct fields and set them in the map.
	//inspecting fileds and set a new map
	for i := 0; i < rVal.NumField(); i++ {
		// get prop | attribute with name and type
		field := rTyp.Field(i)
		filedValue := rVal.Field(i)
		//  convert it to interface{}
		result[field.Name] = filedValue.Interface()
	}
	return result
}

// Deserialize populates a struct from a Result map.
func Deserialize(r Result, refOut interface{}) error {
	//  Ensure refOut is a pointer
	rVal := reflect.ValueOf(refOut)
	if rVal.Kind() != reflect.Pointer {
		return fmt.Errorf("refOut must be a pointer !")
	}

	// Get the type of the struct  ex:Person struct
	structVal := rVal.Elem()
	if structVal.Kind() != reflect.Struct {
		return fmt.Errorf("refOut must be a pointer struct !")

	}

	// Get the type of the struct
	structType := structVal.Type()

	for i := 0; i < structType.NumField(); i++ {
		//  acces struct fileds
		field := structType.Field(i)

		// Check if the field exists in the map
		if val, ok := r[field.Name]; ok {
			structVal.Field(i).Set(reflect.ValueOf(val))
		}

	}
	return nil
}

// Example usage:
type Person struct {
	Name string
	Age  int
}

func Example() {
	// Serialization example
	p := Person{
		Name: "khaled",
		Age:  28,
	}

	serializedRes := Serialize(p)

	// Print serialized result
	fmt.Printf("Serialized result type: %s\n", reflect.TypeOf(serializedRes).Kind())
	for k, v := range serializedRes {
		fmt.Printf("key: %v, value: %v\n", k, v)
	}

	// Deserialization example
	pr := &Person{} // Create an instance and pass its address
	err := Deserialize(serializedRes, pr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the deserialized struct
	fmt.Printf("Deserialized struct: %+v\n", pr)
}
