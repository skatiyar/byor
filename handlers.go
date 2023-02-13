package byor

const (
	RES_OK  = 0 // Status OK
	RES_ERR = 1 // Status Error
	RES_NX  = 2 // Status Not Found
)

func requestHandler(req []byte) (int, string) {
	cmds, parseErr := parseReq(req)
	if parseErr != nil {
		return RES_ERR, parseErr.Error()
	}
	if len(cmds) < 1 {
		return RES_ERR, "Invalid number of commands"
	}
	switch cmds[0] {
	case "ping":
		return RES_OK, "pong"
	default:
		return RES_ERR, "Invalid command"
	}
}

func responseHandler(res []byte) (int, string) {
	return parseRes(res)
}

func validResponseOnComposeError() []byte {
	res := make([]byte, 0)
	res = appendVarint(res, int32(RES_ERR))
	res = append(res, []byte("Server response length invalid")...)
	return res
}
