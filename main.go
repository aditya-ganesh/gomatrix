package main

import (
	"fmt"
	"math/rand"
	"time"
	"unicode/utf8"
)

var (
	firstKana = 0xFF66
	lastKana  = 0xFF9D
)

func randRange(min, max int) int {
	return rand.Intn(max+1-min) + min
}

func main() {
	rand.Seed(time.Now().UnixNano())
	character := randRange(firstKana, lastKana)
	printstring := utf8.AppendRune(nil, rune(character))
	fmt.Printf(string(printstring))
}
