package internal

import "fmt"

func getString(m map[string]any, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("field '%s' is required", key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("field '%s' must be a non-empty string", key)
	}
	return s, nil
}
func getBool(m map[string]any, key string) (bool, error) {
	v, ok := m[key]
	if !ok {
		return false, fmt.Errorf("field '%s' is required", key)
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("field '%s' must be a boolean", key)
	}
	return b, nil
}
func getInt(m map[string]any, key string) (int, error) {
	v, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("field '%s' is required", key)
	}
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("field '%s' must be a number", key)
	}
	i := int(f)
	if i <= 0 {
		return 0, fmt.Errorf("field '%s' must be > 0", key)
	}

	return i, nil
}
