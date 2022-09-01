package screen

import (
	"eklase/state"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/nfnt/resize"
)

func Instruction(th *material.Theme, state *state.State) Screen {
	var (
		close    widget.Clickable
		addimage widget.Image
	)
	addimage.Src = drawPNG()
	return func(gtx layout.Context) (Screen, layout.Dimensions) {
		widgetcolor := th.ContrastBg // to change the background color of the widget
		widgetcolor.A, widgetcolor.R, widgetcolor.G, widgetcolor.B = 0xff, 0x00, 0x0a, 0x12
		max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)
		paint.FillShape(gtx.Ops, widgetcolor, clip.Rect{Max: max}.Op())

		textRowLayout := func(gtx layout.Context) layout.Dimensions {
			labelFn := func(size float32, s string) material.LabelStyle { // to change the color of the text
				style := material.Label(th, unit.Dp(size), s)
				style.Color = color.NRGBA{A: 0xff, R: 0xff, G: 0xff, B: 0xff}
				return style
			}
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
				layout.Rigid(rowInset(labelFn(18, "* INSTRUCTION TO MY 'TETRIS':").Layout)),
				layout.Rigid(rowInset(labelFn(15, "").Layout)),
				layout.Rigid(rowInset(labelFn(15, "- ARROW LEFT (←)    : MOVE THE FIGURE TO THE LEFT").Layout)),
				layout.Rigid(rowInset(labelFn(15, "- ARROW RIGHT (→) : MOVE THE FIGURE TO THE RIGHT").Layout)),
				layout.Rigid(rowInset(labelFn(15, "- ARROW UP (↑)         : ROTATE THE FIGURE").Layout)),
				layout.Flexed(1, rowInset(labelFn(15, "- ARROW DOWN (↓)   : SPEED UP THE FALL OF THE FIGURE").Layout)),
			)
		}

		matCloseBut := material.Button(th, &close, " Close ")
		matCloseBut.Font = text.Font{Variant: "Mono", Weight: text.Bold, Style: text.Italic}

		layout.SE.Layout(gtx, addimage.Layout)
		d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, textRowLayout),
			layout.Rigid(rowInset(matCloseBut.Layout)),
		)
		if close.Clicked() {
			return MainMenu(th, state), d
		}
		return nil, d
	}
}

func drawPNG() paint.ImageOp {
	f, err := os.Open("options.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	m := resize.Resize(100, 100, img, resize.Lanczos3)

	src := paint.NewImageOp(m)
	return src
}
