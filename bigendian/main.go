package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	b := make([]byte, 8)
	var a uint64 = 12

	binary.BigEndian.PutUint64(b, a)
	fmt.Printf("Big: %0X\n", b) // output for debug

	binary.LittleEndian.PutUint64(b, a)
	fmt.Printf("Little: %0X\n", b) // output for debug
}
