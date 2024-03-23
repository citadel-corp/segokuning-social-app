package request

import (
	"net/url"
	"slices"
	"strconv"
)

func CheckPositiveInt(params url.Values, key string) (int, bool) {
	var result int

	if params.Has(key) {
		var value = params.Get(key)

		if value == "" {
			return 0, false
		}

		i, err := strconv.Atoi(value)
		if err != nil {
			return 0, false
		}

		if i < 0 {
			return 0, false
		}

		result = i
	}

	return result, true
}

func CheckBoolean(params url.Values, key string) (bool, bool) {
	var result bool

	if params.Has(key) {
		var value = params.Get(key)

		if value == "" {
			return false, false
		}

		i, err := strconv.ParseBool(value)
		if err != nil {
			return false, false
		}

		result = i
	}

	return result, true
}

func CheckEnum(params url.Values, key string, enum []string) (string, bool) {
	var result string

	if params.Has(key) {
		var value = params.Get(key)

		if value == "" {
			return "", false
		}

		if !slices.Contains(enum, value) {
			return "", false
		}

		result = value
	}

	return result, true
}
