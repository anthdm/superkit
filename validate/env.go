package validate

import (
	"fmt"
	"os"
)

type EnvValidator struct {
	key          string
	defaultValue string
	rules        []RuleSet
}

func Env(key string, rules ...RuleSet) *EnvValidator {
	ruleSets := make([]RuleSet, len(rules))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = rules[i]
	}
	return &EnvValidator{rules: ruleSets, key: key}
}

func (e *EnvValidator) Default(defaultValue string) *EnvValidator {
	e.defaultValue = defaultValue
	return e
}

func (e *EnvValidator) Validate() string {
	var fieldName, fieldValue string
	val := os.Getenv(e.key)

	fieldName = e.key
	fieldValue = val

	for _, set := range e.rules {
		set.FieldValue = fieldValue
		set.FieldName = fieldName
		if !set.ValidateFunc(set) {
			msg := set.MessageFunc(set)
			if len(set.ErrorMessage) > 0 {
				msg = set.ErrorMessage
			}
			if e.defaultValue == "" {
				panic(fmt.Sprintf("Error parsing env %s: %s", e.key, msg))
			} else {
				return e.defaultValue
			}
		}
	}

	if val == "" {
		return e.defaultValue
	}

	return val
}
