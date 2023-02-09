package byor

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

const (
	HEADER_BYTES   = 4
	MAX_CMD_LENGTH = math.MaxInt32
)

func Encoder(wrt io.Writer, cmd string) error {
	dataSize := len(cmd)
	if dataSize < 1 || dataSize > MAX_CMD_LENGTH {
		return errors.New("Invalid cmd length")
	}
	if headErr := binary.Write(wrt, binary.LittleEndian, int32(dataSize)); headErr != nil {
		return headErr
	}
	if cmdErr := binary.Write(wrt, binary.LittleEndian, []byte(cmd)); cmdErr != nil {
		return cmdErr
	}
	return nil
}

func Decoder(rdr io.Reader) (string, error) {
	var dataSize int32
	if dErr := binary.Read(rdr, binary.LittleEndian, &dataSize); dErr != nil {
		return "", dErr
	}
	cmd := make([]byte, dataSize)
	if cErr := binary.Read(rdr, binary.LittleEndian, &cmd); cErr != nil {
		return "", cErr
	}
	return string(cmd), nil
}
