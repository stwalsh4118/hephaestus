package deploy

import "github.com/stwalsh4118/hephaestus/backend/internal/model"

// DiffResult holds the result of comparing two diagram node sets.
type DiffResult struct {
	Added     []model.DiagramNode
	Removed   []model.DiagramNode
	Unchanged []model.DiagramNode
}

// ComputeDiff compares current and incoming diagram nodes by ID and returns
// which nodes were added, removed, or unchanged. Config changes are not
// tracked — a node with the same ID is considered unchanged regardless of
// config.
func ComputeDiff(current, incoming []model.DiagramNode) DiffResult {
	currentMap := make(map[string]model.DiagramNode, len(current))
	for _, n := range current {
		currentMap[n.ID] = n
	}

	incomingMap := make(map[string]model.DiagramNode, len(incoming))
	for _, n := range incoming {
		incomingMap[n.ID] = n
	}

	var result DiffResult

	// Nodes in incoming but not in current are added.
	// Nodes in both are unchanged.
	for id, node := range incomingMap {
		if _, exists := currentMap[id]; exists {
			result.Unchanged = append(result.Unchanged, node)
		} else {
			result.Added = append(result.Added, node)
		}
	}

	// Nodes in current but not in incoming are removed.
	for id, node := range currentMap {
		if _, exists := incomingMap[id]; !exists {
			result.Removed = append(result.Removed, node)
		}
	}

	return result
}
