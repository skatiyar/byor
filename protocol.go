package byor

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

const (
	MAX_BYTES = math.MaxInt32
)

func Encoder(wrt io.Writer, cmd []byte) error {
	dataSize := len(cmd)
	if dataSize < 1 || dataSize > MAX_BYTES {
		return errors.New("invalid cmd length")
	}
	if headErr := binary.Write(wrt, binary.LittleEndian, int32(dataSize)); headErr != nil {
		return headErr
	}
	if cmdErr := binary.Write(wrt, binary.LittleEndian, cmd); cmdErr != nil {
		return cmdErr
	}
	return nil
}

func Decoder(rdr io.Reader) ([]byte, error) {
	var dataSize int32
	if dErr := binary.Read(rdr, binary.LittleEndian, &dataSize); dErr != nil {
		return nil, dErr
	}
	cmd := make([]byte, dataSize)
	if cErr := binary.Read(rdr, binary.LittleEndian, &cmd); cErr != nil {
		return nil, cErr
	}
	return cmd, nil
}
