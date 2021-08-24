package snmp

import (
	"fmt"
	"regexp"
	"strconv"
)

// ResultValue represent a snmp value
type ResultValue struct {
	SubmissionType string      // used when sending the metric
	Value          interface{} // might be a `string` or `float64` type
}

func (sv *ResultValue) toFloat64() (float64, error) {
	switch sv.Value.(type) {
	case float64:
		return sv.Value.(float64), nil
	case string:
		val, err := strconv.ParseFloat(sv.Value.(string), 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse `%s`: %s", sv.Value, err.Error())
		}
		return val, nil
	}
	return 0, fmt.Errorf("invalid type %T for value %#v", sv.Value, sv.Value)
}

func (sv ResultValue) toString() (string, error) {
	switch sv.Value.(type) {
	case float64:
		return strconv.Itoa(int(sv.Value.(float64))), nil
	case string:
		return sv.Value.(string), nil
	}
	return "", fmt.Errorf("invalid type %T for value %#v", sv.Value, sv.Value)
}

func (sv ResultValue) extractStringValue(extractValuePattern *regexp.Regexp) (ResultValue, error) {
	switch sv.Value.(type) {
	case string:
		srcValue := sv.Value.(string)
		matches := extractValuePattern.FindStringSubmatch(srcValue)
		if matches == nil {
			return ResultValue{}, fmt.Errorf("extract value extractValuePattern does not match (extractValuePattern=%v, srcValue=%v)", extractValuePattern, srcValue)
		}
		matchedValue := matches[1] // use first matching group
		return ResultValue{SubmissionType: sv.SubmissionType, Value: matchedValue}, nil
	default:
		return sv, nil
	}
}
