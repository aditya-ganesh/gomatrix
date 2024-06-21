package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

var (
	firstKana = 0xFF66
	lastKana  = 0xFF9D
)

type Raindrop struct {
	text string
	x    int
	y    int
}

func randRange(min, max int) int {
	return rand.Intn(max+1-min) + min
}

func makeRainDrop(maxLength, screenWidth int) Raindrop {

	dropLength := randRange(1, maxLength)
	raindropString := make([]rune, dropLength)

	for i := range raindropString {
		character := randRange(firstKana, lastKana)
		raindropString[i] = rune(character)
	}

	drop := Raindrop{
		string(raindropString),
		rand.Intn(screenWidth),
		-rand.Intn(maxLength),
	}
	return (drop)
}

func drawRaindrop(s tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		row++
	}
}

func main() {

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault
	s.SetStyle(defStyle)

	w, h := s.Size()
	interval := time.Duration(1e6*0.1) * time.Microsecond

	raindropCount := 10
	var raindrops []Raindrop

	maxDropLength := 8

	for range raindropCount {
		raindrop := makeRainDrop(maxDropLength, w)
		raindrops = append(raindrops, raindrop)
	}

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Event handler goroutine

	for {
		s.Clear()
		time.Sleep(interval)

		for i := range raindropCount {

			drawRaindrop(s, raindrops[i].x, raindrops[i].y, defStyle, raindrops[i].text)
		}
		s.Show()

		for i := range len(raindrops) {
			if raindrops[i].y > h {
				raindrops[i].y = -rand.Intn(maxDropLength)
			} else {
				raindrops[i].y++
			}
		}

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}

}
