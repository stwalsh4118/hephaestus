package storage

import (
	"errors"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// ErrNotFound is returned when a diagram ID does not exist in the store.
var ErrNotFound = errors.New("diagram not found")

// DiagramStore defines the persistence operations for diagrams.
type DiagramStore interface {
	// Create persists a new diagram and returns it with a generated ID.
	Create(d *model.Diagram) (*model.Diagram, error)

	// Get retrieves a diagram by ID. Returns ErrNotFound if it does not exist.
	Get(id string) (*model.Diagram, error)

	// Update replaces an existing diagram. Returns ErrNotFound if the ID does not exist.
	Update(id string, d *model.Diagram) (*model.Diagram, error)
}
