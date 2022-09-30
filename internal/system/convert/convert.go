package convert

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// FloatFromString converts a string to a float64. Will return error if unable to convert.
func FloatFromString(raw interface{}) (float64, error) {
	str, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("unable to parse, value not string: %T", raw)
	}
	flt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("could not convert value: %s Error: %s", str, err)
	}
	return flt, nil
}

// IntFromString converts a string to an int. Will return error if unable to convert.
func IntFromString(raw interface{}) (int, error) {
	str, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("unable to parse, value not string: %T", raw)
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("unable to parse as int: %T", raw)
	}
	return n, nil
}

// Int64FromString converts a string to an int64. Will return error if unable to convert.
func Int64FromString(raw interface{}) (int64, error) {
	str, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("unable to parse, value not string: %T", raw)
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse as int64: %T", raw)
	}
	return n, nil
}

// TimeFromUnixTimestampFloat converts a unix timestamp in float form to
// // a time.Time struct. Will return error if unable to convert.
func TimeFromUnixTimestampFloat(raw interface{}) (time.Time, error) {
	ts, ok := raw.(float64)
	if !ok {
		return time.Time{}, fmt.Errorf("unable to parse, value not float64: %T", raw)
	}
	return time.UnixMilli(int64(ts)), nil
}

// TimeFromUnixTimestampDecimal converts a unix timestamp in decimal form to
// a time.Time struct.
func TimeFromUnixTimestampDecimal(input float64) time.Time {
	i, f := math.Modf(input)
	return time.Unix(int64(i), int64(f*(1e9)))
}

// UnixTimestampToTime converts a unix timestamp to a time.Time struct.
func UnixTimestampToTime(timeint64 int64) time.Time {
	return time.Unix(timeint64, 0)
}

// UnixTimestampStrToTime converts a unix timestamp in str format to
// a time.Timestruct. Will return error if unable to convert.
func UnixTimestampStrToTime(timeStr string) (time.Time, error) {
	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0), nil
}

// BoolPtr takes in boolean condition and returns pointer version of it
func BoolPtr(condition bool) *bool {
	b := condition
	return &b
}

// RoundDuration takes in a time.Duration, and rounds to the wanted digits
func RoundDuration(d time.Duration, digits int) time.Duration {
	var divs = []time.Duration{time.Duration(1), time.Duration(10), time.Duration(100), time.Duration(1000)}
	switch {
	case d > time.Second:
		d = d.Round(time.Second / divs[digits])
	case d > time.Millisecond:
		d = d.Round(time.Millisecond / divs[digits])
	case d > time.Microsecond:
		d = d.Round(time.Microsecond / divs[digits])
	}
	return d
}
