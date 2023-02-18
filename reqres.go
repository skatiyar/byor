package byor

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	LEN_BYTES = 4
	// Max int32 size minus 4 to account for strings number prefix
	// as max length of data sent can be MaxInt32
	MAX_STR_LENGTH = math.MaxInt32 - LEN_BYTES
)

const (
	SER_NIL = 0
	SER_ERR = 1
	SER_STR = 2
	SER_INT = 3
	SER_ARR = 4
)

func composeReq(raw string) ([]byte, error) {
	return serializeStringSlice(strings.Split(raw, " "))
}

func parseReq(raw []byte) ([]string, error) {
	return deserializeStringSlice(raw)
}

func composeRes(status int, raw interface{}) ([]byte, error) {
	res := make([]byte, 0)
	res = appendVarint(res, int32(status))

	switch val := raw.(type) {
	case nil:
		res = appendVarint(res, SER_NIL)
	case int32:
		res = appendVarint(res, SER_INT)
		res = appendVarint(res, val)
	case int:
		res = appendVarint(res, SER_INT)
		res = appendVarint(res, int32(val))
	case string:
		if len(val) > MAX_STR_LENGTH {
			return nil, errors.New("invalid length for response: " + val)
		}
		res = appendVarint(res, SER_STR)
		res = append(res, []byte(val)...)
	case []string:
		res = appendVarint(res, SER_ARR)
		sbytes, sErr := serializeStringSlice(val)
		if sErr != nil {
			return nil, sErr
		}
		res = append(res, sbytes...)
	default:
		return nil, errors.New("invalid data type")
	}

	return res, nil
}

func parseRes(raw []byte) (int, []string) {
	if len(raw) < LEN_BYTES {
		return RES_ERR, []string{"invalid server response length"}
	}

	status := int(getVarint(raw[0:4]))
	dtype := getVarint(raw[4:8])
	switch dtype {
	case SER_NIL:
		return status, []string{"<nil>"}
	case SER_INT:
		return status, []string{strconv.Itoa(int(getVarint(raw[8:])))}
	case SER_STR:
		return status, []string{string(raw[8:])}
	case SER_ARR:
		strs, strsErr := deserializeStringSlice(raw[8:])
		if strsErr != nil {
			return status, []string{strsErr.Error()}
		}
		return status, strs
	default:
		return status, []string{"invalid data type received"}
	}
}

func serializeStringSlice(raw []string) ([]byte, error) {
	result := make([]byte, 0)
	numStrs := 0
	for idx := range raw {
		lenCmd := len(raw[idx])
		if lenCmd > 0 {
			result = appendVarint(result, int32(lenCmd))
			result = append(result, []byte(raw[idx])...)
			numStrs += 1
		}
	}
	if numStrs != 0 && len(result) <= MAX_STR_LENGTH {
		numBin := make([]byte, 0)
		numBin = appendVarint(numBin, int32(numStrs))
		return append(numBin, result...), nil
	}
	return nil, errors.New("invalid string slice")
}

func deserializeStringSlice(raw []byte) ([]string, error) {
	if len(raw) < (LEN_BYTES * 2) {
		return nil, errors.New("invalid command")
	}
	numStrs := int(getVarint(raw[0:4]))
	strs := make([]string, 0)
	pointer := 4
	for pointer < len(raw) {
		strLen := int(getVarint(raw[pointer : pointer+4]))
		pointer += 4
		strs = append(strs, string(raw[pointer:pointer+strLen]))
		pointer += strLen
	}
	if numStrs < len(strs) || numStrs > len(strs) {
		return nil, errors.New("invalid expected number of commands")
	}
	return strs, nil
}
