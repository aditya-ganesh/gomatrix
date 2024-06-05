package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	firstKana = 0xFF66
	lastKana  = 0xFF9D
)

func randRange(min, max int) int {
	return rand.Intn(max+1-min) + min
}

func makeRainDrop(minLength, maxLength int) string {
	dropLength := randRange(minLength, maxLength)
	raindrop := make([]rune, dropLength)

	for i := range raindrop {
		character := randRange(firstKana, lastKana)
		raindrop[i] = rune(character)
	}
	return string(raindrop)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	raindrop := makeRainDrop(5, 8)
	fmt.Println(raindrop)
}
