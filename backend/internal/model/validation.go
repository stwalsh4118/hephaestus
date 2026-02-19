package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidationError holds all validation failures for a diagram.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Errors, "; "))
}

// ValidateDiagram checks that a Diagram has all required fields and valid values.
func ValidateDiagram(d *Diagram) error {
	var errs []string

	if d.ID == "" {
		errs = append(errs, "id is required")
	}
	if d.Name == "" {
		errs = append(errs, "name is required")
	}
	if d.Nodes == nil {
		errs = append(errs, "nodes is required")
	}
	if d.Edges == nil {
		errs = append(errs, "edges is required")
	}

	for i, node := range d.Nodes {
		errs = append(errs, validateNode(i, &node)...)
	}

	for i, edge := range d.Edges {
		errs = append(errs, validateEdge(i, &edge)...)
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

func validateNode(index int, n *DiagramNode) []string {
	var errs []string
	prefix := fmt.Sprintf("nodes[%d]", index)

	if n.ID == "" {
		errs = append(errs, fmt.Sprintf("%s.id is required", prefix))
	}
	if n.Type == "" {
		errs = append(errs, fmt.Sprintf("%s.type is required", prefix))
	} else if !ValidServiceTypes[n.Type] {
		errs = append(errs, fmt.Sprintf("%s.type %q is not a valid service type", prefix, n.Type))
	}
	if n.Name == "" {
		errs = append(errs, fmt.Sprintf("%s.name is required", prefix))
	}
	if n.Position == nil {
		errs = append(errs, fmt.Sprintf("%s.position is required", prefix))
	}

	if len(n.Config) > 0 {
		errs = append(errs, validateConfig(prefix, n.Type, n.Config)...)
	}

	return errs
}

func validateEdge(index int, e *DiagramEdge) []string {
	var errs []string
	prefix := fmt.Sprintf("edges[%d]", index)

	if e.ID == "" {
		errs = append(errs, fmt.Sprintf("%s.id is required", prefix))
	}
	if e.Source == "" {
		errs = append(errs, fmt.Sprintf("%s.source is required", prefix))
	}
	if e.Target == "" {
		errs = append(errs, fmt.Sprintf("%s.target is required", prefix))
	}

	return errs
}

func validateConfig(prefix string, nodeType string, raw json.RawMessage) []string {
	var base struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &base); err != nil {
		return []string{fmt.Sprintf("%s.config: invalid JSON: %v", prefix, err)}
	}
	if base.Type != nodeType {
		return []string{fmt.Sprintf("%s.config.type %q does not match node type %q", prefix, base.Type, nodeType)}
	}
	return nil
}
