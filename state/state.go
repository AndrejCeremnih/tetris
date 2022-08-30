package state

import (
	"eklase/storage"
)

// State is the application context (aka state). It provides access to the
// features that do not depend on implementation e.g. (T)UI framework.
type State struct {
	storage *storage.Storage // Provides DB access.

	quit bool // True if the application should exit.
}

// New returns a new state handler. Returns an error if any of the steps fails.
func New(s *storage.Storage) *State {
	return &State{storage: s}
}

func (h *State) DeleteRecordByID(id int) error {
	return h.storage.DeleteRecordByID(id)
}

func (h *State) EditRecordByID(id int, name, surname string) error {
	return h.storage.EditRecordByID(id, name, surname)
}

// Students returns students stored in the database.
func (h *State) Students(name, surname string) ([]storage.StudentEntry, error) {
	return h.storage.Students(name, surname)
}

// AddStudent adds a student to the database.
func (v *State) AddStudent(name, surname string) error {
	return v.storage.AddStudent(name, surname)
}

// Quit requests quitting the application.
func (v *State) Quit() {
	v.quit = true
}

// ShouldQuit returns whether quitting the application was requested.
func (v *State) ShouldQuit() bool {
	return v.quit
}
