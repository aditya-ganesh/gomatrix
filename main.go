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

type raindrop struct {
	text string
	x    int
	y    int
}

func randRange(min, max int) int {
	return rand.Intn(max+1-min) + min
}

func makeRainDrop(maxLength, screenWidth int) raindrop {

	dropLength := randRange(1, maxLength)
	raindropString := make([]rune, dropLength)

	for i := range raindropString {
		character := randRange(firstKana, lastKana)
		raindropString[i] = rune(character)
	}

	drop := raindrop{
		string(raindropString),
		rand.Intn(screenWidth),
		0,
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

	w, _ := s.Size()
	interval := time.Duration(1e6*0.05) * time.Microsecond

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	// Goroutine for refreshing the screen
	refresh := func() {
		time.Sleep(interval)
		s.Show()
		log.Printf("%s", time.Now())
	}

	// Event handler goroutine
	event_handler := func() {
		for {
			// Update screen
			s.Show()

			// Poll event
			ev := s.PollEvent()

			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
				w, _ = s.Size()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					quit()
				}
			}
		}
	}

	// Set default text style
	defStyle := tcell.StyleDefault
	s.SetStyle(defStyle)

	raindrop := makeRainDrop(8, w)
	log.Printf("%d x %d : %s", raindrop.x, raindrop.y, raindrop.text)
	drawRaindrop(s, raindrop.x, raindrop.y, defStyle, raindrop.text)

	raindrop = makeRainDrop(8, w)
	drawRaindrop(s, raindrop.x, raindrop.y, defStyle, raindrop.text)

	// Clear screen

	refresh()
	event_handler()

}
