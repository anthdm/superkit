package validate

import (
	"fmt"
	"reflect"
)

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

func Required() RuleSet {
	return RuleSet{
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
}

func Url() RuleSet {
	return RuleSet{
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
}

func Email() RuleSet {
	return RuleSet{
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
