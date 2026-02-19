package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// Compile-time assertion that FileStore implements DiagramStore.
var _ DiagramStore = (*FileStore)(nil)

const defaultStorageDir = "./data/diagrams"

// ErrInvalidID is returned when a diagram ID contains path traversal characters.
var ErrInvalidID = errors.New("invalid diagram ID")

// FileStore implements DiagramStore using individual JSON files on disk.
type FileStore struct {
	dir string
	mu  sync.RWMutex
}

// NewFileStore creates a FileStore that persists diagrams in the given directory.
// If dir is empty, the default directory is used. The directory is created if it
// does not exist.
func NewFileStore(dir string) (*FileStore, error) {
	if dir == "" {
		dir = defaultStorageDir
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	return &FileStore{dir: dir}, nil
}

// Create persists a new diagram with a generated UUID and returns the stored copy.
func (fs *FileStore) Create(d *model.Diagram) (*model.Diagram, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	out := *d
	out.ID = uuid.New().String()

	if err := fs.writeDiagram(&out); err != nil {
		return nil, fmt.Errorf("create diagram: %w", err)
	}
	return &out, nil
}

// Get retrieves a diagram by ID from disk. Returns ErrNotFound if the file does not exist.
func (fs *FileStore) Get(id string) (*model.Diagram, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	fs.mu.RLock()
	defer fs.mu.RUnlock()

	return fs.readDiagram(id)
}

// Update replaces an existing diagram on disk. Returns ErrNotFound if the ID does not exist.
func (fs *FileStore) Update(id string, d *model.Diagram) (*model.Diagram, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, err := fs.readDiagram(id); err != nil {
		return nil, err
	}

	out := *d
	out.ID = id
	if err := fs.writeDiagram(&out); err != nil {
		return nil, fmt.Errorf("update diagram: %w", err)
	}
	return &out, nil
}

// validateID rejects IDs that could cause path traversal.
func validateID(id string) error {
	if id == "" || filepath.Base(id) != id {
		return ErrInvalidID
	}
	return nil
}

func (fs *FileStore) filePath(id string) string {
	return filepath.Join(fs.dir, id+".json")
}

func (fs *FileStore) readDiagram(id string) (*model.Diagram, error) {
	data, err := os.ReadFile(fs.filePath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("read diagram file: %w", err)
	}

	var d model.Diagram
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("unmarshal diagram: %w", err)
	}
	return &d, nil
}

func (fs *FileStore) writeDiagram(d *model.Diagram) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal diagram: %w", err)
	}

	// Write to temp file first, then rename for atomicity.
	tmp, err := os.CreateTemp(fs.dir, "*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, fs.filePath(d.ID)); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
