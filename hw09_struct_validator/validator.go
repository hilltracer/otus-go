package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	for i, e := range v {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(fmt.Sprintf("%s: %v", e.Field, e.Err))
	}
	return b.String()
}

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrValidation      = errors.New("validation error")
)

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	var all ValidationErrors
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		val := rv.Field(i)

		tag := sf.Tag.Get("validate")
		if tag == "" || !val.CanInterface() {
			continue
		}

		for _, rule := range strings.Split(tag, "|") {
			if err := applyRule(val, rule); err != nil {
				all = append(all, ValidationError{Field: sf.Name, Err: err})
			}
		}
	}

	if len(all) > 0 {
		return fmt.Errorf("%w: %w", ErrValidation, all)
	}
	return nil
}

func applyRule(v reflect.Value, rule string) error {
	kind := v.Kind()
	if kind == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if err := applyRule(v.Index(i), rule); err != nil {
				return fmt.Errorf("elem %d: %w", i, err)
			}
		}
		return nil
	}
	if v.Kind() == reflect.String {
		return checkString(v.String(), rule)
	}
	if v.Kind() == reflect.Int {
		return checkInt(int(v.Int()), rule)
	}
	return nil
}

func checkString(s, rule string) error {
	switch {
	case strings.HasPrefix(rule, "len:"):
		n, _ := strconv.Atoi(rule[4:])
		if len(s) != n {
			return fmt.Errorf("len must be %d", n)
		}
	case strings.HasPrefix(rule, "regexp:"):
		pat := rule[7:]
		re, err := regexp.Compile(pat)
		if err != nil {
			return fmt.Errorf("invalid regexp: %w", err)
		}
		if !re.MatchString(s) {
			return fmt.Errorf("does not match %s", pat)
		}
	case strings.HasPrefix(rule, "in:"):
		for _, v := range strings.Split(rule[3:], ",") {
			if s == v {
				return nil
			}
		}
		return fmt.Errorf("%q is not in set", s)
	}
	return nil
}

func checkInt(n int, rule string) error {
	switch {
	case strings.HasPrefix(rule, "min:"):
		minVal, _ := strconv.Atoi(rule[4:])
		if n < minVal {
			return fmt.Errorf("must be ≥ %d", minVal)
		}
	case strings.HasPrefix(rule, "max:"):
		maxVal, _ := strconv.Atoi(rule[4:])
		if n > maxVal {
			return fmt.Errorf("must be ≤ %d", maxVal)
		}
	case strings.HasPrefix(rule, "in:"):
		for _, v := range strings.Split(rule[3:], ",") {
			x, _ := strconv.Atoi(v)
			if n == x {
				return nil
			}
		}
		return fmt.Errorf("%d is not in set", n)
	}
	return nil
}
