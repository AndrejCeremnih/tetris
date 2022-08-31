package screen

import (
	"eklase/state"
	"image/png"
	"log"
	"os"

	"gioui.org/layout"
	"gioui.org/op/paint"
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
	textRowLayout := func(gtx layout.Context) layout.Dimensions {
		//widgetcolor := th.ContrastBg // To change the widget's background color
		//widgetcolor.A, widgetcolor.R, widgetcolor.G, widgetcolor.B = 0xff, 0x00, 0x0a, 0x12
		//max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)
		//paint.FillShape(gtx.Ops, widgetcolor, clip.Rect{Max: max}.Op())

		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
			layout.Rigid(rowInset(material.Label(th, unit.Dp(18), "* INSTRUCTION TO MY 'TETRIS':").Layout)),
			layout.Rigid(rowInset(material.Label(th, unit.Dp(15), "").Layout)),
			layout.Rigid(rowInset(material.Label(th, unit.Dp(15), "- ARROW LEFT (←)    : MOVE THE FIGURE TO THE LEFT").Layout)),
			layout.Rigid(rowInset(material.Label(th, unit.Dp(15), "- ARROW RIGHT (→) : MOVE THE FIGURE TO THE RIGHT").Layout)),
			layout.Rigid(rowInset(material.Label(th, unit.Dp(15), "- ARROW UP (↑)         : ROTATE THE FIGURE").Layout)),
			layout.Flexed(1, rowInset(material.Label(th, unit.Dp(15), "- ARROW DOWN (↓)   : SPEED UP THE FALL OF THE FIGURE").Layout)),
		)
	}
	buttonRowLayout := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
			layout.Rigid(material.Button(th, &close, "Close").Layout),
		)
	}
	return func(gtx layout.Context) (Screen, layout.Dimensions) {
		layout.SE.Layout(gtx, addimage.Layout)
		d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, textRowLayout),
			layout.Rigid(rowInset(buttonRowLayout)),
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
