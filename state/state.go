package state

// State is the application context (aka state). It provides access to the
// features that do not depend on implementation e.g. (T)UI framework.
type State struct {
	quit bool // True if the application should exit.
}

// New returns a new state handler. Returns an error if any of the steps fails.
func New(s bool) *State {
	return &State{quit: s}
}

// Quit requests quitting the application.
func (v *State) Quit() {
	v.quit = true
}

// ShouldQuit returns whether quitting the application was requested.
func (v *State) ShouldQuit() bool {
	return v.quit
}
