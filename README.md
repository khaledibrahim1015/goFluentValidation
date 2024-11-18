# goFluentValidation
a robust validator package in Go

# Go Struct Validator Package

A flexible and extensible validation package for Go structs that provides both built-in and custom validation rules.

## Table of Contents
- [Installation](#installation)
- [Features](#features)
- [Usage](#usage)
  - [Basic Validation](#basic-validation)
  - [Built-in Validators](#built-in-validators)
  - [Custom Validators](#custom-validators)
- [Validation Rules](#validation-rules)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Installation

```bash
go get github.com/yourusername/validator
```

## Features

- Built-in validators for common use cases
- Custom validator support
- Struct tag-based validation rules
- Clear error messaging
- Support for nested structs
- Type-safe validation
- Chainable validation rules

## Usage

### Basic Validation

1. First, create a new validator instance:

```go
validator := validator.New()
```

2. Define your struct with validation tags:

```go
type User struct {
    Name  string `validate:"required,min=2,max=50"`
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18,max=100"`
}
```

3. Validate your struct:

```go
user := &User{
    Name:  "John",
    Email: "john@example.com",
    Age:   25,
}

if err := validator.Validate(user); err != nil {
    fmt.Println(err)
}
```

### Built-in Validators

The package comes with several built-in validators:

| Validator | Description | Example |
|-----------|-------------|---------|
| `required` | Field cannot be empty | `validate:"required"` |
| `min` | Minimum length for strings or minimum value for numbers | `validate:"min=2"` |
| `max` | Maximum length for strings or maximum value for numbers | `validate:"max=50"` |
| `email` | Valid email format | `validate:"email"` |
| `regex` | Custom regular expression pattern | `validate:"regex=^[0-9]+$"` |

### Custom Validators

You can register custom validators for specific validation logic:

```go
validator := validator.New()

// Register a custom username validator
validator.RegisterCustomValidator("valid_username", func(field reflect.Value) error {
    username := field.String()
    if !strings.Contains(username, "_") {
        return fmt.Errorf("username must contain underscore")
    }
    return nil
})

// Use the custom validator in your struct
type User struct {
    Username string `validate:"required,valid_username"`
}
```

## Validation Rules

### Combining Rules

Multiple validation rules can be combined using commas:

```go
type Person struct {
    Name        string `validate:"required,min=2,max=50"`
    Email       string `validate:"required,email"`
    PhoneNumber string `validate:"required,regex=^01[0125][0-9]{8}$"`
}
```

### Available Rules

- **required**: Field must not be empty or zero value
- **min=X**: 
  - For strings: minimum length
  - For numbers: minimum value
- **max=X**:
  - For strings: maximum length
  - For numbers: maximum value
- **email**: Must be a valid email format
- **regex=pattern**: Must match the specified regular expression pattern

## Error Handling

The validator returns `ValidationErrors` which implements the `error` interface:

```go
type ValidationError struct {
    Field   string
    Message string
}

type ValidationErrors []ValidationError
```

Error messages are formatted as: `"fieldName : errorMessage"`

Example error output:
```
Name : length must be at least 2; Age : value must be at least 18
```

## Examples

### Basic Required Fields

```go
type User struct {
    Name  string `validate:"required"`
    Email string `validate:"required"`
}

func main() {
    v := validator.New()
    usr := &User{}
    if err := v.Validate(usr); err != nil {
        fmt.Println(err)
        // Output: Name : field is required; Email : field is required
    }
}
```

### Advanced Validation

```go
type Person struct {
    Name        string `validate:"required,min=2,max=50"`
    Email       string `validate:"required,email"`
    PhoneNumber string `validate:"required,regex=^01[0125][0-9]{8}$"`
}

func main() {
    v := validator.New()
    
    user := Person{
        Name:        "khaled",
        Email:       "khaled.ibrahem.ahmed.ali@gmail.com",
        PhoneNumber: "01140849506",
    }
    
    if err := v.Validate(&user); err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("Validation passed")
    }
}
```

### Custom Validation Rules

```go
type Department struct {
    Username   string `validate:"required,valid_username"`
    Department string `validate:"required,valid_department"`
    Email      string `validate:"required,email"`
}

func main() {
    v := validator.New()
    
    // Register custom validators
    v.RegisterCustomValidator("valid_username", func(field reflect.Value) error {
        if !strings.Contains(field.String(), "_") {
            return fmt.Errorf("username must contain underscore")
        }
        return nil
    })
    
    v.RegisterCustomValidator("valid_department", func(field reflect.Value) error {
        if !strings.HasPrefix(field.String(), "DEP_") {
            return fmt.Errorf("department must start with DEP_")
        }
        return nil
    })
    
    // Valid data
    dept := Department{
        Username:   "khaled_ibrahim",
        Department: "DEP_tech",
        Email:      "khaled@gmail.com",
    }
    
    if err := v.Validate(&dept); err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("Validation passed")
    }
}
```

## Best Practices

1. Always use pointers when validating structs
2. Keep validation rules simple and composable
3. Use custom validators for complex business logic
4. Handle validation errors appropriately in your application
5. Document custom validation rules clearly
6. Use meaningful field names and error messages

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
