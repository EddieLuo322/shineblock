package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
)

func JsonToArray(jsonString string) []string {
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}

func Byt() {
	stringBytes := []byte("abcdweafwaefwefews")
	fmt.Println(stringBytes)
	fmt.Printf("%x\n", stringBytes)
	// 自动将 byte 数据转化为对应的 字符串
	fmt.Printf("%s\n", stringBytes)
	fmt.Printf("%s\n", string(stringBytes))
	// 将byte 直接转化为hex字符串
	encodeStr := hex.EncodeToString(stringBytes)
	fmt.Printf("%s\n", encodeStr)
	// 将hex(16进制)的字符串直接转化为hex的数据
	decodeBytes, _ := hex.DecodeString(encodeStr)
	fmt.Printf("%x\n", decodeBytes)
}

func main() {
	Byt()
}
