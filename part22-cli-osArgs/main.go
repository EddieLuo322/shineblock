package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Printf("%v\n", args)
	fmt.Printf("%v\n", args[1])
}

// go build -o bc main.go

// bc
// ./bc addBlock -data "send $100 to YaDi"

// ./bc printchain
