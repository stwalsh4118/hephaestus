package model

import (
	"encoding/json"
	"testing"
)

func validDiagram() *Diagram {
	return &Diagram{
		ID:   "test-id",
		Name: "Test Diagram",
		Nodes: []DiagramNode{
			{
				ID:          "node-1",
				Type:        ServiceTypeAPIService,
				Name:        "My API",
				Description: "An API service",
				Position:    &Position{X: 100, Y: 200},
				Config:      json.RawMessage(`{"type":"api-service","endpoints":[],"port":8080}`),
			},
		},
		Edges: []DiagramEdge{
			{
				ID:     "edge-1",
				Source: "node-1",
				Target: "node-2",
				Label:  "connects to",
			},
		},
	}
}

func TestDiagramJSON_RoundTrip(t *testing.T) {
	original := validDiagram()

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Diagram
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID: got %q, want %q", decoded.ID, original.ID)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name: got %q, want %q", decoded.Name, original.Name)
	}
	if len(decoded.Nodes) != len(original.Nodes) {
		t.Fatalf("Nodes length: got %d, want %d", len(decoded.Nodes), len(original.Nodes))
	}

	node := decoded.Nodes[0]
	if node.ID != "node-1" {
		t.Errorf("Node.ID: got %q, want %q", node.ID, "node-1")
	}
	if node.Type != ServiceTypeAPIService {
		t.Errorf("Node.Type: got %q, want %q", node.Type, ServiceTypeAPIService)
	}
	if node.Name != "My API" {
		t.Errorf("Node.Name: got %q, want %q", node.Name, "My API")
	}
	if node.Position == nil {
		t.Fatal("Node.Position is nil")
	}
	if node.Position.X != 100 || node.Position.Y != 200 {
		t.Errorf("Node.Position: got (%f, %f), want (100, 200)", node.Position.X, node.Position.Y)
	}
	if len(node.Config) == 0 {
		t.Error("Node.Config is empty after round-trip")
	}

	edge := decoded.Edges[0]
	if edge.ID != "edge-1" {
		t.Errorf("Edge.ID: got %q, want %q", edge.ID, "edge-1")
	}
	if edge.Source != "node-1" {
		t.Errorf("Edge.Source: got %q, want %q", edge.Source, "node-1")
	}
	if edge.Target != "node-2" {
		t.Errorf("Edge.Target: got %q, want %q", edge.Target, "node-2")
	}
	if edge.Label != "connects to" {
		t.Errorf("Edge.Label: got %q, want %q", edge.Label, "connects to")
	}
}

func TestDiagramJSON_UnmarshalFromFrontendFormat(t *testing.T) {
	input := `{
		"id": "abc-123",
		"name": "Untitled Diagram",
		"nodes": [
			{
				"id": "n1",
				"type": "postgresql",
				"name": "Main DB",
				"description": "Primary database",
				"position": {"x": 50, "y": 75},
				"config": {"type": "postgresql", "engine": "PostgreSQL", "version": "16"}
			}
		],
		"edges": [
			{
				"id": "e1",
				"source": "n1",
				"target": "n2",
				"label": ""
			}
		]
	}`

	var d Diagram
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if d.ID != "abc-123" {
		t.Errorf("ID: got %q, want %q", d.ID, "abc-123")
	}
	if d.Nodes[0].Type != ServiceTypePostgreSQL {
		t.Errorf("Node.Type: got %q, want %q", d.Nodes[0].Type, ServiceTypePostgreSQL)
	}
	if d.Nodes[0].Position.X != 50 || d.Nodes[0].Position.Y != 75 {
		t.Errorf("Position: got (%f, %f), want (50, 75)", d.Nodes[0].Position.X, d.Nodes[0].Position.Y)
	}
}

func TestValidateDiagram_Valid(t *testing.T) {
	d := validDiagram()
	if err := ValidateDiagram(d); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateDiagram_ValidEmptyNodesAndEdges(t *testing.T) {
	d := &Diagram{
		ID:    "test-id",
		Name:  "Test",
		Nodes: []DiagramNode{},
		Edges: []DiagramEdge{},
	}
	if err := ValidateDiagram(d); err != nil {
		t.Errorf("expected no error for empty nodes/edges, got: %v", err)
	}
}

func TestValidateDiagram_MissingID(t *testing.T) {
	d := validDiagram()
	d.ID = ""
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for missing ID")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "id is required")
}

func TestValidateDiagram_MissingName(t *testing.T) {
	d := validDiagram()
	d.Name = ""
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for missing Name")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "name is required")
}

func TestValidateDiagram_NilNodes(t *testing.T) {
	d := validDiagram()
	d.Nodes = nil
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for nil Nodes")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "nodes is required")
}

func TestValidateDiagram_NilEdges(t *testing.T) {
	d := validDiagram()
	d.Edges = nil
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for nil Edges")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "edges is required")
}

func TestValidateDiagram_InvalidServiceType(t *testing.T) {
	d := validDiagram()
	d.Nodes[0].Type = "invalid-service"
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for invalid service type")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, `nodes[0].type "invalid-service" is not a valid service type`)
}

func TestValidateDiagram_EmptyNodeType(t *testing.T) {
	d := validDiagram()
	d.Nodes[0].Type = ""
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for empty node type")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "nodes[0].type is required")
}

func TestValidateDiagram_MissingNodePosition(t *testing.T) {
	d := validDiagram()
	d.Nodes[0].Position = nil
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for missing position")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "nodes[0].position is required")
}

func TestValidateDiagram_MissingNodeName(t *testing.T) {
	d := validDiagram()
	d.Nodes[0].Name = ""
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for missing node name")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "nodes[0].name is required")
}

func TestValidateDiagram_MissingEdgeSource(t *testing.T) {
	d := validDiagram()
	d.Edges[0].Source = ""
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for missing edge source")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "edges[0].source is required")
}

func TestValidateDiagram_MissingEdgeTarget(t *testing.T) {
	d := validDiagram()
	d.Edges[0].Target = ""
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for missing edge target")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, "edges[0].target is required")
}

func TestValidateDiagram_ConfigTypeMismatch(t *testing.T) {
	d := validDiagram()
	d.Nodes[0].Type = ServiceTypeRedis
	d.Nodes[0].Config = json.RawMessage(`{"type":"postgresql"}`)
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected error for config type mismatch")
	}
	ve := err.(*ValidationError)
	assertContains(t, ve.Errors, `nodes[0].config.type "postgresql" does not match node type "redis"`)
}

func TestValidateDiagram_MultipleErrors(t *testing.T) {
	d := &Diagram{}
	err := ValidateDiagram(d)
	if err == nil {
		t.Fatal("expected validation errors")
	}
	ve := err.(*ValidationError)
	if len(ve.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %v", len(ve.Errors), ve.Errors)
	}
}

func TestValidateDiagram_AllServiceTypes(t *testing.T) {
	serviceTypes := []string{
		ServiceTypeAPIService,
		ServiceTypePostgreSQL,
		ServiceTypeRedis,
		ServiceTypeNginx,
		ServiceTypeRabbitMQ,
	}

	for _, st := range serviceTypes {
		t.Run(st, func(t *testing.T) {
			d := &Diagram{
				ID:   "test",
				Name: "test",
				Nodes: []DiagramNode{
					{
						ID:       "n1",
						Type:     st,
						Name:     "test node",
						Position: &Position{X: 0, Y: 0},
					},
				},
				Edges: []DiagramEdge{},
			}
			if err := ValidateDiagram(d); err != nil {
				t.Errorf("service type %q should be valid, got: %v", st, err)
			}
		})
	}
}

func assertContains(t *testing.T, errs []string, target string) {
	t.Helper()
	for _, e := range errs {
		if e == target {
			return
		}
	}
	t.Errorf("expected error %q in %v", target, errs)
}
