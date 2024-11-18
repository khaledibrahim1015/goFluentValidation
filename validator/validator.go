package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	validate    = "validate"
	required    = "required"
	requiredMsg = "field is required"
	min         = "min"
	max         = "max"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string
	Message string
}

// Custom DataTypes
type (
	// ValidationErrors holds multiple validation errors
	ValidationErrors []ValidationError
)

func (ve ValidationErrors) Error() string {

	var errMsgs []string
	for _, errVal := range ve {
		errMsgs = append(errMsgs, fmt.Sprintf("%s : %s", errVal.Field, errVal.Message))
	}
	return strings.Join(errMsgs, "; ")

}

// Validator handles validation logic
type Validator struct {
	errors ValidationErrors
}

// New Create a new Validator instance

func New() *Validator {
	return &Validator{} // errors will be initialized to nil
}

// Validate performs basic validation on the provided struct
func (v *Validator) Validate(s interface{}) error {
	v.errors = ValidationErrors{}

	rVal := reflect.ValueOf(s)
	// Validate type pointer
	if rVal.Kind() != reflect.Pointer {
		return fmt.Errorf("validation requires a struct pointer input")
	}

	// Get the type of the struct  ex:Person struct
	var structVal reflect.Value
	structVal = rVal.Elem()
	// Validate type struct
	if structVal.Kind() != reflect.Struct {
		return fmt.Errorf("refOut must be a pointer struct !")
	}

	// validateFields validates individual fields of the struct
	v.validateFields(structVal)

	if len(v.errors) > 0 {
		return v.errors

	}

	return nil
}

func (v *Validator) validateFields(structVal reflect.Value) {

	// get type
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		currentField := structType.Field(i)
		currentFieldVal := structVal.Field(i)

		// Get validation rules from struct tag `validate:"required,min=2,max=50"`
		tagVal := currentField.Tag.Get(validate)
		if tagVal == "" {
			continue
		}

		rules := strings.Split(tagVal, ",")

		for _, rule := range rules {
			if err := v.applyValidationRule(rule, currentFieldVal, currentField.Name); err != nil {
				v.errors = append(v.errors, ValidationError{
					Field:   currentField.Name,
					Message: err.Error(),
				})
			}
		}

		// emptyStruct  := currentFieldVal.IsZero()
		// if tagVal == required && currentFieldVal.IsZero() {
		// 	v.errors = append(v.errors, ValidationError{
		// 		Field:   currentField.Name,
		// 		Message: requiredMsg,
		// 	})
		// }

	}

}

func (v *Validator) applyValidationRule(rule string, currentFiledVal reflect.Value, fieldName string) error {

	parts := strings.Split(rule, "=")
	ruleName := strings.Trim(parts[0], " ")
	var ruleValue string

	// handle require
	if len(parts) > 1 {
		ruleValue = strings.Trim(parts[1], " ")
	}

	switch ruleName {
	case required:
		if currentFiledVal.IsZero() {
			return fmt.Errorf("field is required")
		}
	case min:
		return v.validateMin(currentFiledVal, ruleValue)
	case max:
		return v.validateMax(currentFiledVal, ruleValue)
	}

	return nil
}

func (v *Validator) validateMin(currentFieldVal reflect.Value, minVlaue string) error {

	min, err := strconv.Atoi(minVlaue)
	if err != nil {
		return fmt.Errorf("invalid min value")
	}

	switch currentFieldVal.Kind() {
	case reflect.String:
		if len(currentFieldVal.String()) < min {
			return fmt.Errorf("length must be at least %d", min)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if currentFieldVal.Int() < int64(min) {
			return fmt.Errorf("value must be at least %d", min)
		}
	}

	return nil

}

func (v *Validator) validateMax(currentFieldVal reflect.Value, maxValue string) error {
	max, err := strconv.Atoi(maxValue)
	if err != nil {
		return fmt.Errorf("invalid max value")
	}

	switch currentFieldVal.Kind() {
	case reflect.String:
		if len(currentFieldVal.String()) > max {
			return fmt.Errorf("length must be at least %d", min)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if currentFieldVal.Int() > int64(max) {
			return fmt.Errorf("value must be at least %d", min)
		}

	}
	return nil

}

type User struct {
	Name  string `validate:"required"`
	Email string `validate:"required"`
}

func ExampleOne() {

	v := New()
	usr := &User{}
	if err := v.Validate(usr); err != nil {
		fmt.Println(err)

	}

}

type Person struct {
	Name string `validate:"required,min=2,max=50"`
	Age  int    `validate:"min=18,max=100"`
}

func ExampleTwo() {

	v := New()
	user := Person{
		Name: "A", // too short
		Age:  15,  // too young
	}

	if err := v.Validate(&user); err != nil {
		fmt.Println(err)
	}
}
