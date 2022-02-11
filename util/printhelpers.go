package util

import (
	"fmt"
	"strings"
)

func PrintHeader(header string) {
	headLength := len([]rune(header)) + 6
	headerBarrier := strings.Repeat("=", headLength)

	fmt.Println(headerBarrier)
	fmt.Printf("!! %s !!\n", header)
	fmt.Println(headerBarrier)
}

func PrintList(stringList []string) {
	for _, s := range stringList {
		fmt.Printf("==> %s\n", s)
	}
}
