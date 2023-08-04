package converter

import (
	"fmt"
	"reflect"
)

// GetMappedFieldValue get the mapped value with converting the original value into target type first
// in case the field mapping does not provide any target type string is used as target type
func (tc typeConverter) GetMappedFieldValue(mapping FieldMapping, originalValue any) (any, error) {
	targetType := "string"
	if mapping.Type != nil {
		targetType = *mapping.Type
	}

	if mapping.FixValue != nil && len(*mapping.FixValue) > 0 {
		return tc.ConvertValue(mapping.FixValue, targetType, tc.convertFunctionMap)
	}

	value := originalValue
	if mapping.ValueMapping != nil && len(mapping.ValueMapping) > 0 {
		for _, valueMapping := range mapping.ValueMapping {
			if valueMapping.Source == originalValue {
				value = valueMapping.Target
			}
		}
	}

	convertedValue, err := tc.ConvertValue(value, targetType, tc.convertFunctionMap)
	if err != nil {
		return nil, err
	}
	if convertedValue == nil {
		return nil, err
	}

	return convertedValue, nil
}

func isInteger(val float64) bool {
	return val == float64(int(val))
}

// ConvertValue convert the input value to input target type by using the given converting functions
// the input converting function map follows the structure => givenType: wishedType: convertingFunction
// Example converting function map for converting string to int and string to float32:
//
//	"string":{
//			"int": convertStringToInt(value),
//			"float32": convertStringToFloat32(value)
//	}
func (tc typeConverter) ConvertValue(value any, targetType string, convertingFunctions map[string]map[string]func(value any) (any, error)) (any, error) {
	if value == nil {
		return nil, nil
	}

	inputValueType := reflect.TypeOf(value).String()
	if inputValueType == targetType {
		return value, nil
	}

	if val, ok := convertingFunctions[inputValueType][targetType]; ok {
		// call function to get the result
		return val(value)
	}

	return nil, fmt.Errorf("cannot map values: unknown type mapping given. inputType: %s, targetType: %s", inputValueType, targetType)
}
