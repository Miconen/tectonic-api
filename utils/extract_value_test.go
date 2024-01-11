package utils

import (
	"reflect"
	"testing"
)

func TestExtractByCloneSuccess(t *testing.T) {
	// Test case 1: Key exists
	inputMap := map[string]string{"key1": "value1", "key2": "value2"}
	expectedValue := "value1"
	expectedMap := map[string]string{"key2": "value2"}

	resultValue, resultMap, err := ExtractByClone(inputMap, "key1")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resultValue != expectedValue {
		t.Errorf("Expected value: %s, got: %s", expectedValue, resultValue)
	}

	if !reflect.DeepEqual(resultMap, expectedMap) {
		t.Errorf("Expected map: %v, got: %v", expectedMap, resultMap)
	}
}

func TestExtractByCloneFail(t *testing.T) {
	// Test case 2: Key does not exist
	inputMap := map[string]string{"key2": "value2"}

	resultValue, resultMap, err := ExtractByClone(inputMap, "key1")

	if err == nil {
		t.Error("Expected an error, but got none")
	}

	if resultValue != "" {
		t.Errorf("Expected empty value, got: %s", resultValue)
	}

	if !reflect.DeepEqual(resultMap, inputMap) {
		t.Errorf("Expected map: %v, got: %v", inputMap, resultMap)
	}
}

func TestExtractByMutateSuccess(t *testing.T) {
	// Test case 1: Key exists
	inputMap := map[string]string{"key1": "value1", "key2": "value2"}
	expectedValue := "value1"

	resultValue, err := ExtractByMutate(inputMap, "key1")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resultValue != expectedValue {
		t.Errorf("Expected value: %s, got: %s", expectedValue, resultValue)
	}

	_, exists := inputMap["key1"]
	if exists {
		t.Error("Expected key1 to be deleted, but it still exists")
	}
}

func TestExtractByMutateFail(t *testing.T) {
	// Test case 2: Key does not exist
	inputMap := map[string]string{"key2": "value2"}

	resultValue, err := ExtractByMutate(inputMap, "key1")

	if err == nil {
		t.Error("Expected an error, but got none")
	}

	if resultValue != "" {
		t.Errorf("Expected empty value, got: %s", resultValue)
	}
}
