package main

import (
	"log"
	"math/rand"
	"os"
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
	len  int
}

func randRange(min, max int) int {
	return rand.Intn(max+1-min) + min
}

func makeRainDrop(minLength, maxLength, screenWidth int) Raindrop {

	dropLength := randRange(minLength, maxLength)
	raindropString := make([]rune, dropLength)

	for i := range raindropString {
		character := randRange(firstKana, lastKana)
		raindropString[i] = rune(character)
	}

	x := rand.Intn(screenWidth)

	drop := Raindrop{
		string(raindropString),
		x,
		-randRange(2*dropLength, 10*dropLength),
		dropLength,
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

	raindropCount := w / 4
	var raindrops []Raindrop

	maxDropLength := h / 2
	minDropLength := h / 5

	for range raindropCount {
		raindrop := makeRainDrop(minDropLength, maxDropLength, w)
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
		os.Exit(0)
	}
	defer quit()

	// Display handler goroutine
	go func() {
		for {
			s.Clear()
			time.Sleep(interval)

			for i := range raindropCount {

				drawRaindrop(s, raindrops[i].x, raindrops[i].y, defStyle, raindrops[i].text)
			}
			s.Show()

			for i := range len(raindrops) {
				// If a raindrop escapes the scene, spawn a new one in its place.
				if raindrops[i].y > h {
					raindrop := makeRainDrop(minDropLength, maxDropLength, w)
					raindrops[i] = raindrop
					// Otherwise just let it drop down the screen
				} else {
					raindrops[i].y++
				}
			}
		}
	}()

	// Event handler goroutine
	eventChan := make(chan tcell.Event)
	go func() {
		for {
			event := s.PollEvent()
			eventChan <- event
		}
	}()

	for {
		select {
		case event := <-eventChan:
			switch ev := event.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlZ, tcell.KeyCtrlC:
					quit()
				}
			}

		}
	}

}
