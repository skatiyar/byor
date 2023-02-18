package byor

import (
	"strings"
)

const (
	RES_OK  = 0 // Status OK
	RES_ERR = 1 // Status Error
	RES_NX  = 2 // Status Not Found
)

var safeMap = NewHashMap(128)

func requestHandler(req []byte) (int, interface{}) {
	cmds, parseErr := parseReq(req)
	if parseErr != nil {
		return RES_ERR, parseErr.Error()
	}
	if len(cmds) < 1 {
		return RES_ERR, "invalid number of commands"
	}
	switch strings.ToLower(cmds[0]) {
	case "ping":
		return RES_OK, "pong"
	case "get":
		if len(cmds) == 2 {
			value := safeMap.Get(cmds[1])
			if value == "" {
				return RES_NX, "key doesn't exists"
			} else {
				return RES_OK, value
			}
		} else {
			return RES_ERR, "invalid parameters for command: get"
		}
	case "set":
		if len(cmds) == 3 {
			safeMap.Put(cmds[1], cmds[2])
			return RES_OK, "key, value saved"
		} else {
			return RES_ERR, "invalid parameters for command: set"
		}
	case "delete":
		if len(cmds) == 2 {
			safeMap.Delete(cmds[1])
			return RES_OK, "key deleted"
		} else {
			return RES_ERR, "invalid parameters for command: del"
		}
	case "size":
		return RES_OK, safeMap.Size()
	case "keys":
		return RES_OK, safeMap.Keys()
	case "purge":
		safeMap.Clear()
		return RES_OK, "hashmap purged"
	default:
		return RES_ERR, "invalid command"
	}
}

func responseHandler(res []byte) (int, []string) {
	return parseRes(res)
}

func validResponseOnComposeError() []byte {
	res := make([]byte, 0)
	res = appendVarint(res, int32(RES_ERR))
	res = append(res, []byte("Server response length invalid")...)
	return res
}
