package byor

import (
	"encoding/binary"
	"errors"
	"math"
	"strings"
)

const (
	LEN_BYTES = 4
	// Max int32 size minus 4 to account for strings number prefix
	// as max length of data sent can be MaxInt32
	MAX_STR_LENGTH = math.MaxInt32 - LEN_BYTES
)

func composeReq(raw string) ([]byte, error) {
	strs := strings.Split(raw, " ")
	result := make([]byte, 0)
	numCmds := 0
	for idx := range strs {
		cmd := strings.Trim(strs[idx], " ")
		lenCmd := len(cmd)
		if lenCmd >= 1 {
			result = binary.LittleEndian.AppendUint32(result, uint32(lenCmd))
			result = append(result, []byte(cmd)...)
			numCmds += 1
		}
	}
	if numCmds != 0 && len(result) <= MAX_STR_LENGTH {
		numBin := make([]byte, 0)
		numBin = binary.LittleEndian.AppendUint32(numBin, uint32(numCmds))
		return append(numBin, result...), nil
	}

	return nil, errors.New("Invalid command")
}

func parseReq(raw []byte) ([]string, error) {
	if len(raw) < (LEN_BYTES * 2) {
		return nil, errors.New("Invalid command")
	}
	numCmds := int(binary.LittleEndian.Uint32(raw[0:4]))
	cmds := make([]string, 0)
	pointer := 4
	for pointer <= len(raw) {
		cmdLen := int(binary.LittleEndian.Uint32(raw[pointer : pointer+4]))
		pointer += 4
		cmds = append(cmds, string(raw[pointer:pointer+cmdLen]))
		pointer += cmdLen
	}
	if numCmds < len(cmds) || numCmds > len(cmds) {
		return nil, errors.New("Invalid expected number of commands")
	}
	return cmds, nil
}

func composeRes(status int, raw string) ([]byte, error) {
	if len(raw) > MAX_STR_LENGTH {
		return nil, errors.New("Invalid length for response: " + raw)
	}
	res := make([]byte, 0)
	res = binary.LittleEndian.AppendUint32(res, uint32(status))
	res = append(res, []byte(raw)...)
	return res, nil
}

func parseRes(raw []byte) (int, string) {
}
