package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	urlRegex   = regexp.MustCompile(`^(https?:\/\/)?(www\.)?([a-zA-Z0-9\-]+\.)+[a-zA-Z]{2,}(\/[a-zA-Z0-9\-._~:\/?#\[\]@!$&'()*+,;=]*)?$`)
)

// RuleSet holds the state of a single rule.
type RuleSet struct {
	Name         string
	RuleValue    any
	FieldValue   any
	FieldName    any
	ErrorMessage string
	MessageFunc  func(RuleSet) string
	ValidateFunc func(RuleSet) bool
}

// Message overrides the default message of a RuleSet
func (set RuleSet) Message(msg string) RuleSet {
	set.ErrorMessage = msg
	return set
}

type Numeric interface {
	int | float64
}

func In[T any](values []T) RuleSet {
	return RuleSet{
		Name:      "in",
		RuleValue: values,
		ValidateFunc: func(set RuleSet) bool {
			for _, value := range values {
				v := set.FieldValue.(T)
				if reflect.DeepEqual(v, value) {
					return true
				}
			}
			return false
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("should be in %v", values)
		},
	}
}

var ContainsUpper = RuleSet{
	Name: "containsUpper",
	ValidateFunc: func(rule RuleSet) bool {
		str, ok := rule.FieldValue.(string)
		if !ok {
			return false
		}
		for _, ch := range str {
			if unicode.IsUpper(rune(ch)) {
				return true
			}
		}
		return false
	},
	MessageFunc: func(set RuleSet) string {
		return "must contain at least 1 uppercase character"
	},
}

var ContainsDigit = RuleSet{
	Name: "containsDigit",
	ValidateFunc: func(rule RuleSet) bool {
		str, ok := rule.FieldValue.(string)
		if !ok {
			return false
		}
		return hasDigit(str)
	},
	MessageFunc: func(set RuleSet) string {
		return "must contain at least 1 numeric character"
	},
}

var ContainsSpecial = RuleSet{
	Name: "containsSpecial",
	ValidateFunc: func(rule RuleSet) bool {
		str, ok := rule.FieldValue.(string)
		if !ok {
			return false
		}
		return hasSpecialChar(str)
	},
	MessageFunc: func(set RuleSet) string {
		return "must contain at least 1 special character"
	},
}

var Required = RuleSet{
	Name: "required",
	MessageFunc: func(set RuleSet) string {
		return "is a required field"
	},
	ValidateFunc: func(rule RuleSet) bool {
		str, ok := rule.FieldValue.(string)
		if !ok {
			return false
		}
		return len(str) > 0
	},
}

var URL = RuleSet{
	Name: "url",
	MessageFunc: func(set RuleSet) string {
		return "is not a valid url"
	},
	ValidateFunc: func(set RuleSet) bool {
		u, ok := set.FieldValue.(string)
		if !ok {
			return false
		}
		return urlRegex.MatchString(u)
	},
}

var Email = RuleSet{
	Name: "email",
	MessageFunc: func(set RuleSet) string {
		return "is not a valid email address"
	},
	ValidateFunc: func(set RuleSet) bool {
		email, ok := set.FieldValue.(string)
		if !ok {
			return false
		}
		return emailRegex.MatchString(email)
	},
}

var Time = RuleSet{
	Name: "time",
	ValidateFunc: func(set RuleSet) bool {
		t, ok := set.FieldValue.(time.Time)
		if !ok {
			return false
		}
		return t.After(time.Time{})
	},
	MessageFunc: func(set RuleSet) string {
		return "is not a valid time"
	},
}

func TimeAfter(t time.Time) RuleSet {
	return RuleSet{
		Name: "timeAfter",
		ValidateFunc: func(set RuleSet) bool {
			t, ok := set.FieldValue.(time.Time)
			if !ok {
				return false
			}
			return t.After(t)
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("is not after %v", set.FieldValue)
		},
	}
}

func TimeBefore(t time.Time) RuleSet {
	return RuleSet{
		Name: "timeBefore",
		ValidateFunc: func(set RuleSet) bool {
			t, ok := set.FieldValue.(time.Time)
			if !ok {
				return false
			}
			return t.Before(t)
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("is not before %v", set.FieldValue)
		},
	}
}

func EQ[T comparable](v T) RuleSet {
	return RuleSet{
		Name:      "eq",
		RuleValue: v,
		ValidateFunc: func(set RuleSet) bool {
			return set.FieldValue.(T) == v
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("should be equal to %v", v)
		},
	}
}

func LTE[T Numeric](n T) RuleSet {
	return RuleSet{
		Name:      "lte",
		RuleValue: n,
		ValidateFunc: func(set RuleSet) bool {
			return set.FieldValue.(T) <= n
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("should be lesser or equal than %v", n)
		},
	}
}

func GTE[T Numeric](n T) RuleSet {
	return RuleSet{
		Name:      "gte",
		RuleValue: n,
		ValidateFunc: func(set RuleSet) bool {
			return set.FieldValue.(T) >= n
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("should be greater or equal than %v", n)
		},
	}
}

func LT[T Numeric](n T) RuleSet {
	return RuleSet{
		Name:      "lt",
		RuleValue: n,
		ValidateFunc: func(set RuleSet) bool {
			return set.FieldValue.(T) < n
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("should be lesser than %v", n)
		},
	}
}

func GT[T Numeric](n T) RuleSet {
	return RuleSet{
		Name:      "gt",
		RuleValue: n,
		ValidateFunc: func(set RuleSet) bool {
			return set.FieldValue.(T) > n
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("should be greater than %v", n)
		},
	}
}

func Max(n int) RuleSet {
	return RuleSet{
		Name:      "max",
		RuleValue: n,
		ValidateFunc: func(set RuleSet) bool {
			str, ok := set.FieldValue.(string)
			if !ok {
				return false
			}
			return len(str) <= n
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("should be maximum %d characters long", n)
		},
	}
}

func Min(n int) RuleSet {
	return RuleSet{
		Name:      "min",
		RuleValue: n,
		ValidateFunc: func(set RuleSet) bool {
			str, ok := set.FieldValue.(string)
			if !ok {
				return false
			}
			return len(str) >= n
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("should be at least %d characters long", n)
		},
	}
}

func hasDigit(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func hasSpecialChar(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return true
		}
	}
	return false
}
