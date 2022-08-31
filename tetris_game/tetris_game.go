package tetris_game

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/nsf/termbox-go"

	_ "modernc.org/sqlite"
)

const (
	figureBody     = '*'
	figureFgColor0 = termbox.ColorGreen
	figureFgColor1 = termbox.ColorLightCyan
	figureFgColor2 = termbox.ColorLightYellow
	// Use the default background color for the figure.
	figureBgColor = termbox.ColorDefault
)

func getFigureFgColor(i int) termbox.Attribute { // to make the color of the figure random (i is a random number)
	if i == 0 {
		return figureFgColor0
	} else if i == 1 {
		return figureFgColor1
	} else if i == 2 {
		return figureFgColor2
	}
	return termbox.Attribute(0)
}

// writeText writes a string to the buffer.
func writeText(x, y int, s string, fg, bg termbox.Attribute) {
	for i, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

// coord is a coordinate on a plane.
type coord struct {
	x, y int
}

// figure is a struct with fields representing a figure.
type figure struct {
	// Position of a figure.
	pos []coord
}

// game represents a state of the game.
type game struct {
	sn figure
	v  coord
	// Game field dimensions.
	fieldWidth, fieldHeight int
}

func getAllTheSides(g game) (leftBorder, rightBorder, up, down int) {
	leftBorder, rightBorder = 0, g.fieldWidth-1
	up, down = 0, g.fieldHeight-1
	if rightBorder > 36 {
		leftBorder, rightBorder = g.fieldWidth/2-17, g.fieldWidth/2+17
	}
	if down > 28 {
		up, down = g.fieldHeight/2-15, g.fieldHeight/2+15
	}
	return
}

// newfigure returns a new struct instance representing a figure.
// The figure is placed in a random position in the game field.
// The movement direction is right.
func newfigure(g game) figure {
	x := g.fieldWidth / 2
	_, _, y, _ := getAllTheSides(g)
	y = y - 1
	g.sn.pos = []coord{{x, y}, {x, y - 1}, {x, y - 2}}
	return figure{g.sn.pos}
}

// newGame returns a new game state.
func newGame() game {
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		// sn:       newfigure(w, h),
		// ap:       newApple(w, h),
		v: coord{0, 1},
	}
}

func hitsTheFloor(g game) bool {
	_, _, _, down := getAllTheSides(g)
	return g.sn.pos[0].y == down-1
}

// drawfigurePosition draws the current figure position (as a debugging
// information) in the buffer.
func drawfigurePosition(g game, i int) {
	str := fmt.Sprintf("(%d, %d)", g.sn.pos[0].x, g.sn.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, getFigureFgColor(i), figureBgColor)
}

//
//func drawScore(g game) {
//	str := fmt.Sprintf("Score: %d", g.ap.score)
//	if g.ap.score < 5 { // the color will change depending on the score
//		writeText((g.fieldWidth-len(str))/2, 0, str, termbox.ColorWhite, figureBgColor)
//	} else if g.ap.score >= 5 && g.ap.score < 10 {
//		writeText((g.fieldWidth-len(str))/2, 0, str, termbox.ColorLightCyan, figureBgColor)
//	} else if g.ap.score >= 10 && g.ap.score < 20 {
//		writeText((g.fieldWidth-len(str))/2, 0, str, termbox.ColorLightGreen, figureBgColor)
//	} else if g.ap.score >= 20 {
//		writeText((g.fieldWidth-len(str))/2, 0, str, termbox.ColorLightMagenta, figureBgColor)
//	}
//}
//

// drawfigure draws the figure in the buffer.
func drawfigure(g game, sn figure, i int) {
	_, _, up, _ := getAllTheSides(g)
	for cnt, pos := range sn.pos {
		if g.sn.pos[cnt].y > up { // the figure will not appear outside the borders
			termbox.SetCell(pos.x, pos.y, figureBody, getFigureFgColor(i), figureBgColor)
		}
	}
}

func drawBorders(g game) {
	leftBorder, rightBorder, up, down := getAllTheSides(g)

	for i := leftBorder; i < rightBorder; i++ {
		termbox.SetBg(i, up, termbox.ColorLightRed)
		termbox.SetBg(i, down, termbox.ColorLightCyan)
	}
	for j := up; j < down; j++ {
		termbox.SetBg(leftBorder, j, termbox.ColorLightMagenta)
		termbox.SetBg(rightBorder, j, termbox.ColorLightGreen)
	}
	termbox.SetBg(leftBorder, up, termbox.ColorLightYellow)
	termbox.SetBg(leftBorder, down, termbox.ColorLightYellow)
	termbox.SetBg(rightBorder, up, termbox.ColorLightYellow)
	termbox.SetBg(rightBorder, down, termbox.ColorLightYellow)
}

// Redraws the terminal.
func draw(g game, i int) {
	// Clear the old "frame".
	termbox.Clear(getFigureFgColor(i), figureBgColor)
	drawfigurePosition(g, i)

	///     drawScore(g)   !!!

	drawfigure(g, g.sn, i)
	drawBorders(g)
	// Update the "frame".
	termbox.Flush()
}

// Makes a move for a figure. Returns a figure with an updated position.
func movefigure(s figure, v coord, fw, fh int) figure {
	copy(s.pos[1:], s.pos[:])
	s.pos[0] = coord{s.pos[0].x, s.pos[0].y + v.y}
	return s
}

func step(g game, i int) game {
	g.sn = movefigure(g.sn, g.v, g.fieldWidth, g.fieldHeight)

	//if g.sn.pos[0] == g.ap.pos {
	//	g.ap.score++
	//	g.sn.pos = append([]coord{{g.sn.pos[0].x, g.sn.pos[0].y}}, g.sn.pos...)
	//}

	draw(g, i)

	if hitsTheFloor(g) { // CHANGE LATER !!!
		//termbox.SetChar(g.sn.pos[0].x, g.sn.pos[0].y, figureBody)
		termbox.SetCell(64, 2, '%', getFigureFgColor(i), figureBgColor)
		g.sn = newfigure(g)
	}
	return g
}

func moveLeft(g game) game {
	leftBorder, _, _, _ := getAllTheSides(g)
	if !reflect.DeepEqual(g.sn.pos[0], coord{leftBorder + 1, g.sn.pos[0].y}) {
		g.sn.pos[0] = coord{g.sn.pos[0].x - 1, g.sn.pos[0].y}
	}
	return g
}

func moveRight(g game) game {
	_, rightBorder, _, _ := getAllTheSides(g)
	if !reflect.DeepEqual(g.sn.pos[0], coord{rightBorder - 1, g.sn.pos[0].y}) {
		g.sn.pos[0] = coord{g.sn.pos[0].x + 1, g.sn.pos[0].y}
	}
	return g
}

func StartTheGame() {
	// Initialize termbox.
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()

	// Other initialization.
	rand.Seed(time.Now().UnixNano())
	g := newGame()
	g.sn = newfigure(g)

	i := rand.Intn(3) // to make the color of the figure random

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(450 * time.Millisecond)
	defer ticker.Stop()

	AlreadyActive := false
	// This is the main event loop.
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowDown:
					if !AlreadyActive {
						ticker = time.NewTicker(90 * time.Millisecond)
						AlreadyActive = true
					} else {
						ticker = time.NewTicker(450 * time.Millisecond)
						AlreadyActive = false
					}
				case termbox.KeyArrowLeft:
					g = moveLeft(g)
				case termbox.KeyArrowRight:
					g = moveRight(g)
				case termbox.KeyEsc:
					return
				}
			}
		case <-ticker.C:
			g = step(g, i)
		}
	}
}
