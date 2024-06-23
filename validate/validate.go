package validate

import (
	"fmt"
	"maps"
	"net/http"
	"reflect"
	"strconv"
	"unicode"
)

// Errors is a map holding all the possible errors that may
// occur during validation.
type Errors map[string][]string

// Any return true if there is any error.
func (e Errors) Any() bool {
	return len(e) > 0
}

// Add adds an error for a specific field
func (e Errors) Add(field string, msg string) {
	if _, ok := e[field]; !ok {
		e[field] = []string{}
	}
	e[field] = append(e[field], msg)
}

// Get returns all the errors for the given field.
func (e Errors) Get(field string) []string {
	return e[field]
}

// Has returns true whether the given field has any errors.
func (e Errors) Has(field string) bool {
	return len(e[field]) > 0
}

// Schema represents a validation schema.
type Schema map[string][]RuleSet

// Merge merges the two given schemas, returning a new Schema.
func Merge(schema, other Schema) Schema {
	newSchema := Schema{}
	maps.Copy(newSchema, schema)
	maps.Copy(newSchema, other)
	return newSchema
}

// Rules is a function that takes any amount of RuleSets
func Rules(rules ...RuleSet) []RuleSet {
	ruleSets := make([]RuleSet, len(rules))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = rules[i]
	}
	return ruleSets
}

// Validate validates data based on the given Schema.
func Validate(data any, fields Schema) (Errors, bool) {
	errors := Errors{}
	return validate(data, fields, errors)
}

// Request parses an http.Request into data and validates it based
// on the given schema.
func Request(r *http.Request, data any, schema Schema) (Errors, bool) {
	errors := Errors{}
	if err := parseRequest(r, data); err != nil {
		errors["_error"] = []string{err.Error()}
	}
	return validate(data, schema, errors)
}

func validate(data any, schema Schema, errors Errors) (Errors, bool) {
	ok := true
	for fieldName, ruleSets := range schema {
		// Uppercase the field name so we never check un-exported fields.
		// But we need to watch out for member fields that are uppercased by
		// the user. For example (URL, ID, ...)
		if !isUppercase(fieldName) {
			fieldName = string(unicode.ToUpper(rune(fieldName[0]))) + fieldName[1:]
		}

		fieldValue := getFieldAndTagByName(data, fieldName)
		for _, set := range ruleSets {
			set.FieldValue = fieldValue
			set.FieldName = fieldName
			fieldName = string(unicode.ToLower([]rune(fieldName)[0])) + fieldName[1:]
			if !set.ValidateFunc(set) {
				ok = false
				msg := set.MessageFunc(set)
				if len(set.ErrorMessage) > 0 {
					msg = set.ErrorMessage
				}
				if _, ok := errors[fieldName]; !ok {
					errors[fieldName] = []string{}
				}
				errors[fieldName] = append(errors[fieldName], msg)
			}
		}
	}
	return errors, ok
}

func getFieldAndTagByName(v any, name string) any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil
	}
	return fieldVal.Interface()
}

func parseRequest(r *http.Request, v any) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form: %v", err)
		}
		val := reflect.ValueOf(v).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			formTag := field.Tag.Get("form")
			formValue := r.FormValue(formTag)

			if formValue == "" {
				continue
			}

			fieldVal := val.Field(i)
			switch fieldVal.Kind() {
			case reflect.Bool:
				// There are cases where frontend libraries use "on" as the bool value
				// think about toggles. Hence, let's try this first.
				if formValue == "on" {
					fieldVal.SetBool(true)
				} else if formValue == "off" {
					fieldVal.SetBool(false)
					return nil
				} else {
					boolVal, err := strconv.ParseBool(formValue)
					if err != nil {
						return fmt.Errorf("failed to parse bool: %v", err)
					}
					fieldVal.SetBool(boolVal)
				}
			case reflect.String:
				fieldVal.SetString(formValue)
			case reflect.Int, reflect.Int32, reflect.Int64:
				intVal, err := strconv.Atoi(formValue)
				if err != nil {
					return fmt.Errorf("failed to parse int: %v", err)
				}
				fieldVal.SetInt(int64(intVal))
			case reflect.Uint, reflect.Uint32, reflect.Uint64:
				intVal, err := strconv.Atoi(formValue)
				if err != nil {
					return fmt.Errorf("failed to parse int: %v", err)
				}
				fieldVal.SetUint(uint64(intVal))
			case reflect.Float64:
				floatVal, err := strconv.ParseFloat(formValue, 64)
				if err != nil {
					return fmt.Errorf("failed to parse float: %v", err)
				}
				fieldVal.SetFloat(floatVal)
			default:
				return fmt.Errorf("unsupported kind %s", fieldVal.Kind())
			}
		}

	}
	return nil
}

func isUppercase(s string) bool {
	for _, ch := range s {
		if !unicode.IsUpper(rune(ch)) {
			return false
		}
	}
	return true
}
