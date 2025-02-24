// Copyright © 2022 Dell Inc. or its subsidiaries. All Rights Reserved.

// Package utils to provide supporting APIs
package utils

func ToPointer[T any](s T) *T {
	return &s
}

func GetPtrValue[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

// PtrOrNilIfEmpty returns a pointer if the value is non-empty, or nil if empty.
// For types with string as the underlying type.
func PtrOrNilIfEmpty[T ~string](value T) *T {
	if value == "" {
		return nil
	}
	return &value
}

func EqualPtrValues[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func AddElementToMap[T any](m *map[string]T, key string, value T) {
	if *m == nil {
		*m = make(map[string]T)
	}
	(*m)[key] = value
}

func Contains[T any](m map[string]T, key string) bool {
	_, ok := m[key]
	return ok
}

func Keys[T any](m map[string]T) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func ExtractMap(item *map[string]interface{}, key string) map[string]interface{} {
	if item == nil {
		return nil
	}
	if value, ok := (*item)[key].(map[string]interface{}); ok {
		return value
	}
	return nil
}

func ExtractString(item *map[string]interface{}, key string) string {
	if item != nil {
		if value, ok := (*item)[key].(string); ok {
			return value
		}
	}
	return ""
}
