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
	figureFgColor3 = termbox.ColorLightMagenta
	figureFgColor4 = termbox.ColorLightRed
	// Use the default background color for the figure.
	figureBgColor = termbox.ColorDefault
)

var (
	shouldCreateNewFigure    bool
	shouldDrawPreviousFigure bool
	lastFigurePosition       []coord
	lastFigureFgColor        termbox.Attribute
)

func getFigureFgColor(i int) termbox.Attribute { // to make the color of the figure random (i is a random number)
	if i == 0 {
		return figureFgColor0
	} else if i == 1 {
		return figureFgColor1
	} else if i == 2 {
		return figureFgColor2
	} else if i == 3 {
		return figureFgColor3
	} else if i == 4 {
		return figureFgColor4
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
	pos   []coord
	score int
}

// game represents a state of the game.
type game struct {
	fg figure
	v  coord
	// Game field dimensions.
	fieldWidth, fieldHeight int
}

func getAllTheSides(g game) (leftBorder, rightBorder, up, down int) {
	leftBorder, rightBorder = 0, g.fieldWidth-1
	up, down = 1, g.fieldHeight-1 // 'up' is at least 1 to make room for 'Score'
	if rightBorder > 36 {
		leftBorder, rightBorder = g.fieldWidth/2-17, g.fieldWidth/2+17
	}
	if down > 28 {
		up, down = g.fieldHeight/2-15, g.fieldHeight/2+15
	}
	return
}

// newfigure returns a new struct instance representing a figure.
// The figure is placed in the center of the game field from above.
// The movement direction is down.
func newfigure(g game) figure {
	x := g.fieldWidth / 2
	_, _, y, _ := getAllTheSides(g)
	y = y - 1
	g.fg.pos = []coord{{x, y}, {x, y - 1}, {x, y - 2}}
	return figure{g.fg.pos, g.fg.score}
}

// newGame returns a new game state.
func newGame() game {
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		v:           coord{0, 1},
	}
}

func hitsTheFloor(g game) (bool, bool) {
	_, _, _, down := getAllTheSides(g)
	if g.fg.pos[0].y == down-1 {
		shouldCreateNewFigure, shouldDrawPreviousFigure = true, true
	}
	return shouldCreateNewFigure, shouldDrawPreviousFigure
}

// drawfigurePosition draws the current figure position (as a debugging
// information) in the buffer.
func drawfigurePosition(g game, i int) {
	str := fmt.Sprintf("(%d, %d)", g.fg.pos[0].x, g.fg.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, getFigureFgColor(i), figureBgColor)
}

func drawScore(g game) {
	_, _, up, _ := getAllTheSides(g)
	str := fmt.Sprintf("Score: %d", 0)
	if g.fg.score < 5 { // the color will change depending on the score
		writeText((g.fieldWidth-len(str))/2, up-1, str, termbox.ColorWhite, figureBgColor)
	} else if g.fg.score >= 5 && g.fg.score < 10 {
		writeText((g.fieldWidth-len(str))/2, up-1, str, termbox.ColorLightCyan, figureBgColor)
	} else if g.fg.score >= 10 && g.fg.score < 20 {
		writeText((g.fieldWidth-len(str))/2, up-1, str, termbox.ColorLightGreen, figureBgColor)
	} else if g.fg.score >= 20 {
		writeText((g.fieldWidth-len(str))/2, up-1, str, termbox.ColorLightMagenta, figureBgColor)
	}
}

// drawfigure draws the figure in the buffer.
func drawfigure(g game, fg figure, i int) {
	_, _, up, _ := getAllTheSides(g)
	for cnt, pos := range fg.pos {
		if g.fg.pos[cnt].y > up { // the figure will not appear outside the borders
			termbox.SetChar(pos.x, pos.y, figureBody)
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

func drawPreviousFigure(lastFigurePosition []coord, lastFigureColor termbox.Attribute) { // only works for the last figure
	if shouldDrawPreviousFigure {
		for _, pos := range lastFigurePosition {
			termbox.SetCell(pos.x, pos.y, figureBody, lastFigureFgColor, figureBgColor)
		}
	}
}

// Redraws the terminal.
func draw(g game, i int) {
	// Clear the old "frame".
	termbox.Clear(getFigureFgColor(i), figureBgColor)
	drawfigurePosition(g, i)
	drawScore(g)
	drawfigure(g, g.fg, i)
	drawBorders(g)

	hitsTheFloor(g)
	if shouldCreateNewFigure {
		lastFigurePosition, lastFigureFgColor = g.fg.pos, getFigureFgColor(i)
	}
	drawPreviousFigure(lastFigurePosition, lastFigureFgColor)

	// Update the "frame".
	termbox.Flush()
}

// Makes a move for a figure. Returns a figure with an updated position.
func moveFigure(s figure, v coord, fw, fh int) figure {
	copy(s.pos[1:], s.pos[:])
	s.pos[0] = coord{s.pos[0].x, s.pos[0].y + v.y}
	return s
}

func step(g game, i int) (game, int) {
	g.fg = moveFigure(g.fg, g.v, g.fieldWidth, g.fieldHeight)

	draw(g, i)

	hitsTheFloor(g)
	if shouldCreateNewFigure { // CHANGE LATER !!!
		rand.Seed(time.Now().UnixNano())
		ii := rand.Intn(4) // to change the color of the figure
		if i == ii {
			i = i + 1
		} else {
			i = ii
		}
		g.fg = newfigure(g)
		shouldCreateNewFigure = false
	}
	return g, i
}

func moveLeft(g game) game {
	leftBorder, _, _, _ := getAllTheSides(g)
	if !reflect.DeepEqual(g.fg.pos[0], coord{leftBorder + 1, g.fg.pos[0].y}) {
		g.fg.pos[0] = coord{g.fg.pos[0].x - 1, g.fg.pos[0].y}
	}
	return g
}

func moveRight(g game) game {
	_, rightBorder, _, _ := getAllTheSides(g)
	if !reflect.DeepEqual(g.fg.pos[0], coord{rightBorder - 1, g.fg.pos[0].y}) {
		g.fg.pos[0] = coord{g.fg.pos[0].x + 1, g.fg.pos[0].y}
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
	g.fg = newfigure(g)

	i := rand.Intn(5) // to make the color of the figure random

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
			g, i = step(g, i)
		}
	}
}
