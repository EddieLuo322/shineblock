package main

import "fmt"

func main() {
	//block := BLC.NewBlock("Genenis Blocks", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	//	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	//fmt.Println(block)
	s := []int{1, 2, 3, 4}
	s1 := s
	s2 := s[:]
	fmt.Println(s1, s2)
	s[0] = 100
	fmt.Println(s, s1, s2)
	s1[1] = 200
	fmt.Println(s, s1, s2)
	s2[2] = 300
	fmt.Println(s, s1, s2)
}
