package utils

import "strings"

func StringSliceToInterfaceMap(input []string) map[string]interface{} {
	const pair = 2
	m := make(map[string]interface{})
	for _, config := range input {
		parts := strings.Split(config, "=")
		if len(parts) >= pair {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func StringSliceToStringMap(input []string) map[string]string {
	const pair = 2
	m := make(map[string]string)
	for _, config := range input {
		parts := strings.Split(config, "=")
		if len(parts) >= pair {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
