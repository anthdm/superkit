package validate

import (
	"fmt"
	"maps"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	urlRegex   = regexp.MustCompile(`^(http(s)?://)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*/?$`)
)

type RuleFunc func() RuleSet

type RuleSet struct {
	Name         string
	RuleValue    any
	FieldValue   any
	FieldName    any
	ErrorMessage string
	MessageFunc  func(RuleSet) string
	ValidateFunc func(RuleSet) bool
}

func (set RuleSet) Message(msg string) RuleSet {
	set.ErrorMessage = msg
	return set
}

type Errors map[string][]string

func (e Errors) Any() bool {
	return len(e) > 0
}

func (e Errors) Add(name string, msg string) {
	if _, ok := e[name]; !ok {
		e[name] = []string{}
	}
	e[name] = append(e[name], msg)
}

func (e Errors) Get(name string) []string {
	return e[name]
}

func (e Errors) Has(name string) bool {
	return len(e[name]) > 0
}

type Schema map[string][]RuleSet

func (schema Schema) Merge(other Schema) Schema {
	newSchema := Schema{}
	maps.Copy(newSchema, schema)
	maps.Copy(newSchema, other)
	return newSchema
}

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
		// reflect panics on un-exported variables.
		if !unicode.IsUpper(rune(fieldName[0])) {
			errors[fieldName] = []string{
				"cant marshal unexported field",
			}
			return errors, false
		}
		fieldValue := getFieldValueByName(data, fieldName)
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

func getFieldValueByName(v any, name string) any {
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
			case reflect.Int:
				intVal, err := strconv.Atoi(formValue)
				if err != nil {
					return fmt.Errorf("failed to parse int: %v", err)
				}
				fieldVal.SetInt(int64(intVal))
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
