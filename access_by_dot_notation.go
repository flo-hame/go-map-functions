package mapFunctions

import (
	"regexp"
	"strconv"
	"strings"
)

// GetValueByFieldPathDotNotation accessing a field value within a map by using dot notation
// It is also possible to access fields within nested slices.
// Beware of case sensitivity
// Examples:
// Given Map: product: { name: "abc", prices: [{"id": 1, "value": 10}, {"id": 2, "value": 15}], related: {product: {"name": "def"}}
// Accessing the products name by fieldPath: product.name
// Accessing the products related product name by fieldPath: product.related.product.name
// Accessing the products second price by fieldPath: product.prices.[1].value
func GetValueByFieldPathDotNotation(fieldPath string, mapStructure map[string]any) (any, error) {
	paths := strings.Split(fieldPath, ".")
	var value = mapStructure[paths[0]]

	for i, path := range strings.Split(fieldPath, ".") {
		if i == 0 {
			continue
		}

		indexAccess, err := getIndexAccessNumber(path)
		if err != nil {
			return nil, err
		}
		if indexAccess != nil {
			switch typedValue := value.(type) {
			case []any:
				value = typedValue[*indexAccess]
			case []map[string]any:
				value = typedValue[*indexAccess]
			case []string:
				value = typedValue[*indexAccess]
			case []int:
				value = typedValue[*indexAccess]
			case []int32:
				value = typedValue[*indexAccess]
			case []int64:
				value = typedValue[*indexAccess]
			case []float32:
				value = typedValue[*indexAccess]
			case []float64:
				value = typedValue[*indexAccess]
			}
			continue
		}

		switch value.(type) {
		case map[string]any:
			value = value.(map[string]any)[path]
		case []map[string]any:
			return getListMap(path, value.([]map[string]any))
		}
	}

	return value, nil
}

func getListMap(path string, inputValues []map[string]any) (any, error) {
	var values []any
	for _, ele := range inputValues {
		value, err := GetValueByFieldPathDotNotation(path, ele)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func getIndexAccessNumber(path string) (*int, error) {
	var re = regexp.MustCompile(`(?m)\[\d\]`)
	matches := re.FindAllString(path, -1)
	if len(matches) == 1 {
		// have to be one match eg. [0] or [1]

		// get only the number between the brackets
		re = regexp.MustCompile(`(?m)\d`)
		matches = re.FindAllString(path, -1)
		index, err := strconv.Atoi(matches[0])
		if err != nil {
			return nil, err
		}

		return &index, nil
	}

	return nil, nil
}
