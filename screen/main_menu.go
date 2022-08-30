package screen

import (
	"eklase/state"
	"eklase/tetris_game"

	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/nfnt/resize"
	"github.com/nsf/termbox-go"
)

// MainMenu defines a main menu screen layout.
func MainMenu(th *material.Theme, state *state.State) Screen {
	var (
		start       widget.Clickable
		instruction widget.Clickable
		quit        widget.Clickable
		addimage    widget.Image
	)
	addimage.Src = draw()
	return func(gtx layout.Context) (Screen, layout.Dimensions) {
		widgetcolor := th.ContrastBg // To change the widget's background color
		widgetcolor.A, widgetcolor.R, widgetcolor.G, widgetcolor.B = 0xff, 0x00, 0x0a, 0x12
		max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)
		paint.FillShape(gtx.Ops, widgetcolor, clip.Rect{Max: max}.Op())

		matStartBut := material.Button(th, &start, "Start")
		matStartBut.Font = text.Font{Variant: "Mono", Weight: text.Bold, Style: text.Italic}
		matStartBut.Background = color.NRGBA{A: 0xff, R: 0x2e, G: 0x7d, B: 0x32}
		matInstructionBut := material.Button(th, &instruction, "Instruction")
		matInstructionBut.Font = text.Font{Variant: "Mono", Weight: text.Bold, Style: text.Italic}
		matQuitBut := material.Button(th, &quit, "Quit")
		matQuitBut.Font = text.Font{Variant: "Smallcaps", Style: text.Italic}
		matQuitBut.Background = color.NRGBA{A: 0xff, R: 0xc6, G: 0x28, B: 0x28}

		d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(rowInset(matStartBut.Layout)),
			layout.Rigid(rowInset(matInstructionBut.Layout)),
			layout.Rigid(rowInset(matQuitBut.Layout)),
		)
		layout.SE.Layout(gtx, addimage.Layout)
		if start.Clicked() {
			state.Quit() // I don't know why, but it doesn't close the widget immediately
			tetris_game.StartTheGame()
		}
		if instruction.Clicked() {
			return Instruction(th, state), d
		}
		if quit.Clicked() {
			state.Quit()
		}
		return nil, d
	}
}

func draw() paint.ImageOp {
	f, err := os.Open("tetris.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// width, height := 100, 10
	width, height := termbox.Size()
	m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	src := paint.NewImageOp(m)
	return src
}
