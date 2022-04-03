package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

func ConvertToHex(value interface{}) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, value)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
