package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validate          = "validate"
	required          = "required"
	requiredMsg       = "field is required"
	min               = "min"
	max               = "max"
	email             = "email"
	regex             = "regex"
	emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
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

	// CustomValidatorFunc is a type for custom validation functions
	CustomValidatorFunc func(field reflect.Value) error
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
	errors           ValidationErrors
	customValidators map[string]CustomValidatorFunc
}

// New Create a new Validator instance

func New() *Validator {
	return &Validator{
		customValidators: make(map[string]CustomValidatorFunc),
	}
}

// RegisterCustomValidator registers a custom validation function
func (v *Validator) RegisterCustomValidator(tagVal string, fn CustomValidatorFunc) {
	v.customValidators[tagVal] = fn
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

	// Second pass: apply custom validators
	v.applyCustomValidators(structVal)

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
	case email:
		if !v.isMatchedRegex(currentFiledVal.String(), emailRegexPattern) {
			return fmt.Errorf("invalid email format")
		}
	case regex:
		if !v.isMatchedRegex(currentFiledVal.String(), ruleValue) {
			return fmt.Errorf("value does not match required format")
		}
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

func (v *Validator) isMatchedRegex(value, pattern string) bool {

	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func (v *Validator) applyCustomValidators(structVal reflect.Value) {

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
			// Check if this rule is a custom validator
			if validator, ok := v.customValidators[rule]; ok {
				// execute validator
				if err := validator(currentFieldVal); err != nil {
					v.errors = append(v.errors, ValidationError{
						Field:   currentField.Name,
						Message: err.Error(),
					})
				}
			}

		}

	}

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

func ExampleThree() {
	type Khaled struct {
		Name        string `validate:"required,min=2,max=50"`
		Email       string `validate:"required,email"`
		PhoneNumber string `validate:"required,regex=^01[0125][0-9]{8}$"`
	}

	v := New()

	// invalid
	user := Khaled{
		Name:        "John",
		Email:       "invalid-email", // invalid format
		PhoneNumber: "123456789",     // invalid format
	}

	if err := v.Validate(&user); err != nil {
		fmt.Println(err)
	}

	// valid
	user2 := Khaled{
		Name:        "khaled",
		Email:       "khaled.ibrahem.ahmed.ali@gmail.com",
		PhoneNumber: "01140849506",
	}
	if err := v.Validate(&user2); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("passed")
	}

}

func ExampleFour() {

	type Test struct {
		Username   string `validate:"required,valid_username"`   // Descriptive tag name
		Department string `validate:"required,valid_department"` // Different validation
		Email      string `validate:"required,email"`
	}

	v := New()
	v.RegisterCustomValidator("valid_username", func(field reflect.Value) error {
		usernameValue := field.String()
		if !strings.Contains(usernameValue, "_") {
			fmt.Errorf("username must contain underscore")
		}
		return nil
	})

	v.RegisterCustomValidator("valid_department", func(field reflect.Value) error {
		deptValue := field.String()
		if !strings.HasPrefix(deptValue, "DEP_") {
			return fmt.Errorf("department must start with DEPT_")
		}
		return nil

	})

	// invalid
	tst := Test{
		Username:   "khaledibrahim", // invalid: missing underscore
		Department: "tech",          // invalid: missing DEPT_ prefix
		Email:      "khaled@mail",
	}
	if err := v.Validate(&tst); err != nil {
		fmt.Println(err)
	}

	// valid
	tst2 := Test{
		Username:   "khaled_ibrahim", // invalid: missing underscore
		Department: "DEP_tech",       // invalid: missing DEPT_ prefix
		Email:      "khaled@gmail.com",
	}
	if err := v.Validate(&tst2); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("passed ")
	}

}
