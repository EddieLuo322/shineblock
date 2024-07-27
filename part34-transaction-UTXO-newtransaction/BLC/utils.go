package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

func Int64ToBytes(num int64) []byte {
	//numStr := strconv.FormatInt(num, 10)
	//numBytes := []byte(numStr)
	//return numBytes
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func JsonToArray(jsonString string) []string {
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}
