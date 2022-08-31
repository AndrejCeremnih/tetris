package main

import (
	"log"
	"os"

	"eklase/screen"
	"eklase/state"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func main() {
	// Run the main event loop.
	go func() {
		w := app.NewWindow(app.Title("Main Menu"))
		if err := mainLoop(w); err != nil {
			log.Fatalf("failed to handle events: %v", err)
		}
		// Gracefully exit the application at the end.
		os.Exit(0)
	}()
	app.Main()
}

func mainLoop(w *app.Window) error {
	var s bool
	appState := state.New(s)

	th := material.NewTheme(gofont.Collection())
	currentLayout := screen.MainMenu(th, appState)

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&op.Ops{}, e)
				layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					nextLayout, d := currentLayout(gtx)
					if nextLayout != nil {
						currentLayout = nextLayout
					}
					return d
				})
				if appState.ShouldQuit() {
					w.Perform(system.ActionClose)
				}
				e.Frame(gtx.Ops)
			case system.DestroyEvent:
				return e.Err
			}
		}
	}
}
