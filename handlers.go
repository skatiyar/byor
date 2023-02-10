package byor

import "encoding/binary"

const (
	RES_OK  = 0 // Status OK
	RES_ERR = 1 // Status Error
	RES_NX  = 2 // Status Not Found
)

func requestHandler(req []byte) (int, string) {
	cmds, parseErr := parseReq(req)
	if parseErr != nil {
		slog.errorf("While parsing commands %v", parseErr)
		return RES_ERR, parseErr.Error()
	}
	if len(cmds) < 1 {
		return RES_ERR, ""
	}
}

func responseHandler(res []byte) (int, string) {

}

func validResponseOnComposeError() []byte {
	res := make([]byte, 0)
	res = binary.LittleEndian.AppendUint32(res, uint32(RES_ERR))
	res = append(res, []byte("Server response length invalid")...)
	return res
}
