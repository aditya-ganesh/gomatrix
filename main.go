package main

import (
	"flag"
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

func drawRaindrop(s tcell.Screen, x, y int, text string, style tcell.Style) {
	row := y
	col := x
	style = style.Dim(true)
	for i, r := range text {

		// TODO: Find a less horrid way of doing this
		if i >= (8) {
			style = style.Dim(false)
		}
		if i >= (len(text) - 4) {
			style = style.Bold(true)
		}
		s.SetContent(col, row, r, nil, style)
		row++
	}
}

func main() {

	// Process command line flags
	var maxLen = flag.Float64("max", 0.5, "Maximum drop length")
	var minLen = flag.Float64("min", 0.2, "Minimum drop length")
	var density = flag.Float64("d", 0.75, "Raindrop density")
	var refresh = flag.Float64("r", 0.1, "Refresh interval in seconds")
	var colour = flag.String("c", "", "Raindrop colour")

	flag.Parse()

	// Start a new TCell screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault

	if *colour != "" {
		tcellColour := tcell.GetColor(*colour)
		defStyle = defStyle.Foreground(tcellColour)
	}
	s.SetStyle(defStyle)

	// Get the screen size to calculate raindrop properties
	w, h := s.Size()

	// Define a refresh interval for screen painting
	interval := time.Duration(1e6**refresh) * time.Microsecond

	// TODO : Make these parameterizable

	maxDropLength := int(float64(h) * *maxLen)
	minDropLength := int(float64(h) * *minLen)

	raindropCount := int(float64(w) * *density)
	var raindrops []Raindrop

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

				drawRaindrop(s, raindrops[i].x, raindrops[i].y, raindrops[i].text, defStyle)
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
