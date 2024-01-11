package utils

// ConvertStringMapToInterfaceMap converts a map[string]string to map[string]interface{}
func StringMapToInterfaceMap(input map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range input {
		result[key] = value
	}

	return result
}
