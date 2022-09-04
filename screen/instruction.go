package screen

import (
	"eklase/state"
	"fmt"
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

func InstructionText(th *material.Theme, list *widget.List) func(gtx layout.Context) layout.Dimensions {
	lightContrast := th.ContrastBg
	lightContrast = color.NRGBA{A: 0xaa, R: 0x38, G: 0x38, B: 0x34}
	darkContrast := th.ContrastBg
	darkContrast = color.NRGBA{A: 0xaa, R: 0x2d, G: 0x2d, B: 0x2d}

	texts1 := []string{"- ARROW LEFT (←)", "- ARROW RIGHT (→)", "- ARROW UP (↑)", "- ARROW DOWN (↓)"}
	texts2 := []string{": MOVE THE FIGURE TO THE LEFT", ": MOVE THE FIGURE TO THE RIGHT", ": ROTATE THE FIGURE", ": SPEED UP THE FALL OF THE FIGURE"}

	return func(gtx layout.Context) layout.Dimensions {
		return material.List(th, list).Layout(gtx, 4, func(gtx layout.Context, index int) layout.Dimensions {
			labelFn := func(size float32, s string) material.LabelStyle { // to change the color of the text
				style := material.Label(th, unit.Dp(size), s)
				style.Color = color.NRGBA{A: 0xff, R: 0xff, G: 0xff, B: 0xff}
				return style
			}

			text1, text2 := texts1[index], texts2[index]

			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					color := lightContrast
					if index%2 == 0 {
						color = darkContrast
					}

					max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)
					paint.FillShape(gtx.Ops, color, clip.Rect{Max: max}.Op())
					return layout.Dimensions{Size: gtx.Constraints.Min}
				}),
				layout.Stacked(rowInset(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(rowInset(labelFn(16.5, text1).Layout)),
					)
				})),
				layout.Stacked(rowInset(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(rowInset(material.Body1(th, fmt.Sprintln("                                     ")).Layout)),
						layout.Rigid(rowInset(labelFn(16.5, text2).Layout)),
					)
				})),
			)
		})
	}
}

func Instruction(th *material.Theme, state *state.State) Screen {
	var (
		close    widget.Clickable
		addimage widget.Image
	)
	addimage.Src = drawPNG()
	list := widget.List{List: layout.List{Axis: layout.Vertical}}

	textRowLayout := InstructionText(th, &list)
	return func(gtx layout.Context) (Screen, layout.Dimensions) {
		widgetcolor := th.ContrastBg // to change the background color of the widget
		widgetcolor.A, widgetcolor.R, widgetcolor.G, widgetcolor.B = 0xff, 0x00, 0x0a, 0x12
		max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)
		paint.FillShape(gtx.Ops, widgetcolor, clip.Rect{Max: max}.Op())

		labelFn := func(size float32, s string) material.LabelStyle { // to change the color of the text
			style := material.Label(th, unit.Dp(size), s)
			style.Color = color.NRGBA{A: 0xff, R: 0xff, G: 0xff, B: 0xff}
			return style
		}

		matCloseBut := material.Button(th, &close, " Close ")
		matCloseBut.Font = text.Font{Variant: "Mono", Weight: text.Bold, Style: text.Italic}

		d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(rowInset(labelFn(21.5, " * INSTRUCTION TO MY 'TETRIS' :").Layout)),
			layout.Rigid(rowInset(labelFn(16.5, "").Layout)),
			layout.Flexed(1, rowInset(textRowLayout)),
			layout.Rigid(rowInset(matCloseBut.Layout)),
		)
		layout.NE.Layout(gtx, addimage.Layout) // NE (North East) is used to move the picture to the top right corner
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
