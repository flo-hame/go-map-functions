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

	convertedValue, err := tc.ConvertValue(originalValue, targetType, tc.convertFunctionMap)
	if err != nil {
		return nil, err
	}
	if convertedValue == nil && mapping.Default != nil {
		return tc.ConvertValue(*mapping.Default, targetType, tc.convertFunctionMap)
	} else if convertedValue == nil {
		return nil, nil
	}

	if mapping.ValueMapping != nil && len(mapping.ValueMapping.Mapping) > 0 {
		if reflect.TypeOf(convertedValue).Kind() == reflect.Slice {
			return tc.convertSliceValueMapping(convertedValue, mapping.ValueMapping.Mapping, targetType)
		}

		for _, valueMapping := range mapping.ValueMapping.Mapping {
			source, err := tc.ConvertValue(valueMapping.Source, targetType, tc.convertFunctionMap)
			if err != nil {
				return nil, err
			}
			if source == convertedValue {
				return tc.ConvertValue(valueMapping.Target, targetType, tc.convertFunctionMap)
			}
		}
		if mapping.ValueMapping.Default != nil {
			return tc.ConvertValue(*mapping.ValueMapping.Default, targetType, tc.convertFunctionMap)
		}
		return nil, nil
	}

	return convertedValue, nil
}

func (tc typeConverter) convertSliceValueMapping(convertedValue any, valueMappings []Mapping, targetType string) (any, error) {
	var convertedValues []any
	convertedVal := reflect.ValueOf(convertedValue)
	for i := 0; i < convertedVal.Len(); i++ {
		for _, valueMapping := range valueMappings {
			source, err := tc.ConvertValue(valueMapping.Source, targetType, tc.convertFunctionMap)
			if err != nil {
				return nil, err
			}
			v := source
			if reflect.TypeOf(source).Kind() == reflect.Slice {
				v = reflect.ValueOf(source).Index(0).Interface()
			}
			if v == convertedVal.Index(i).Interface() {
				convertedValues = append(convertedValues, source)
			}
		}
	}
	return convertedValues, nil
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
