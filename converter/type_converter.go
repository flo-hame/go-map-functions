package converter

import (
	"strconv"
	"strings"
	"time"
)

type TypeConverter interface {
	GetMappedFieldValue(mapping FieldMapping, originalValue any) (any, error)
	ConvertValue(value any, targetType string, convertingFunctions map[string]map[string]func(value any) (any, error)) (any, error)

	GetStringFromAny(value any) (any, error)
	GetStringPtrFromString(value any) (any, error)
	GetStringFromStringPtr(value any) (any, error)
	ConvertDefaultDateTimeStringToTime(value any) (any, error)
	ConvertStringToInt(value any) (any, error)
	ConvertInt64ToString(value any) (any, error)
	ConvertIntToString(value any) (any, error)
	ConvertUint8SliceToString(value any) (any, error)
	ConvertUint8SliceToFloat64(value any) (any, error)
	ConvertTimeToString(value any) (any, error)
	ConvertStringToStringSlice(value any) (any, error)
	ConvertFloat64ToFloat32(value any) (any, error)
	ConvertFloat64ToInt64(value any) (any, error)
	ConvertStringPtrToInt(value any) (any, error)
}

type typeConverter struct {
	convertFunctionMap map[string]map[string]func(value any) (any, error)
}

func NewTypeConverter(convertFunctionMap map[string]map[string]func(value any) (any, error)) TypeConverter {
	tc := typeConverter{
		convertFunctionMap: convertFunctionMap,
	}

	if convertFunctionMap == nil {
		// input to target type map
		defaultFunctionMap := map[string]map[string]func(value any) (any, error){
			"string": {
				"varchar":                tc.GetStringFromAny,
				"*string":                tc.GetStringPtrFromString,
				"int":                    tc.ConvertStringToInt,
				"datetime":               tc.ConvertDefaultDateTimeStringToTime,
				"americanDateTimeString": tc.ConvertInternationalDateTimeStringToAmerican,
				"[]string":               tc.ConvertStringToStringSlice,
			},
			"*string": {
				"string":  tc.GetStringFromStringPtr,
				"varchar": tc.GetStringFromStringPtr,
				"int":     tc.ConvertStringPtrToInt,
			},
			"int": {
				"string":  tc.ConvertIntToString,
				"varchar": tc.ConvertIntToString,
			},
			"int64": {
				"string": tc.ConvertInt64ToString,
			},
			"[]uint8": {
				"string": tc.ConvertUint8SliceToString,
			},
			"time.Time": {
				"string": tc.ConvertTimeToString,
			},
			"float64": {
				"float32": tc.ConvertFloat64ToFloat32,
				"int64":   tc.ConvertFloat64ToInt64,
				"int":     tc.ConvertFloat64ToInt,
			},
		}
		tc.convertFunctionMap = defaultFunctionMap
	}

	return tc
}

func (typeConverter) GetStringFromAny(value any) (any, error) {
	return value, nil
}

func (typeConverter) GetStringPtrFromString(value any) (any, error) {
	return &value, nil
}

func (typeConverter) GetStringFromStringPtr(value any) (any, error) {
	v := value.(*string)
	if v == nil {
		return nil, nil
	}
	return *v, nil
}

func (typeConverter) ConvertDefaultDateTimeStringToTime(value any) (any, error) {
	return time.Parse("2006-01-02 15:04:05", value.(string))
}

func (typeConverter) ConvertInternationalDateTimeStringToAmerican(value any) (any, error) {
	convertedTime, err := time.Parse("2006-01-02 15:04:05", value.(string))
	if err != nil {
		return nil, err
	}

	return convertedTime.Format("2006-02-01 15:04:05"), nil
}

func (typeConverter) ConvertStringToInt(value any) (any, error) {
	return strconv.Atoi(value.(string))
}

func (typeConverter) ConvertInt64ToString(value any) (any, error) {
	return strconv.Itoa(int(value.(int64))), nil
}

func (typeConverter) ConvertIntToString(value any) (any, error) {
	return strconv.Itoa(value.(int)), nil
}

func (typeConverter) ConvertUint8SliceToString(value any) (any, error) {
	return string(value.([]uint8)), nil
}

func (typeConverter) ConvertUint8SliceToFloat64(value any) (any, error) {
	return strconv.ParseFloat(string(value.([]uint8)), 64)
}

func (typeConverter) ConvertTimeToString(value any) (any, error) {
	return value.(time.Time).String(), nil
}

func (typeConverter) ConvertStringToStringSlice(value any) (any, error) {
	return strings.Split(value.(string), ","), nil
}

func (typeConverter) ConvertFloat64ToFloat32(value any) (any, error) {
	return float32(value.(float64)), nil
}

func (typeConverter) ConvertFloat64ToInt64(value any) (any, error) {
	return int64(value.(float64)), nil
}

func (typeConverter) ConvertFloat64ToInt(value any) (any, error) {
	return int(value.(float64)), nil
}

func (tc typeConverter) ConvertStringPtrToInt(value any) (any, error) {
	v := value.(*string)
	if v == nil {
		return nil, nil
	}
	return tc.ConvertStringToInt(*v)
}
