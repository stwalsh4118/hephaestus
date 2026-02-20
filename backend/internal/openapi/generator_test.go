package openapi

import (
	"encoding/json"
	"testing"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

func TestGenerateSpec_SingleGETEndpoint(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/users", ResponseSchema: `{"type":"array","items":{"type":"object"}}`},
	}

	data, err := GenerateSpec(endpoints, "User Service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc spec
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if doc.OpenAPI != openAPIVersion {
		t.Errorf("expected openapi %q, got %q", openAPIVersion, doc.OpenAPI)
	}
	if doc.Info.Title != "User Service" {
		t.Errorf("expected title %q, got %q", "User Service", doc.Info.Title)
	}
	if len(doc.Paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(doc.Paths))
	}
	pi, ok := doc.Paths["/users"]
	if !ok {
		t.Fatal("expected path /users")
	}
	if _, ok := pi["get"]; !ok {
		t.Error("expected GET operation on /users")
	}
}

func TestGenerateSpec_MultipleMethodsOnSamePath(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/users", ResponseSchema: `{"type":"array"}`},
		{Method: "POST", Path: "/users", ResponseSchema: `{"type":"object"}`},
	}

	data, err := GenerateSpec(endpoints, "Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc spec
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(doc.Paths) != 1 {
		t.Fatalf("expected 1 path (grouped), got %d", len(doc.Paths))
	}
	pi := doc.Paths["/users"]
	if _, ok := pi["get"]; !ok {
		t.Error("expected GET operation")
	}
	if _, ok := pi["post"]; !ok {
		t.Error("expected POST operation")
	}
}

func TestGenerateSpec_MultipleDifferentPaths(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/users", ResponseSchema: `{"type":"array"}`},
		{Method: "GET", Path: "/orders", ResponseSchema: `{"type":"array"}`},
	}

	data, err := GenerateSpec(endpoints, "Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc spec
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(doc.Paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(doc.Paths))
	}
	if _, ok := doc.Paths["/users"]; !ok {
		t.Error("expected path /users")
	}
	if _, ok := doc.Paths["/orders"]; !ok {
		t.Error("expected path /orders")
	}
}

func TestGenerateSpec_ValidJSONResponseSchema(t *testing.T) {
	schema := `{"type":"object","properties":{"id":{"type":"integer"},"name":{"type":"string"}}}`
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/item", ResponseSchema: schema},
	}

	data, err := GenerateSpec(endpoints, "Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc spec
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	op := doc.Paths["/item"]["get"]
	resp := op.Responses["200"]
	mt := resp.Content[contentTypeJSON]

	var parsed map[string]any
	if err := json.Unmarshal(mt.Schema, &parsed); err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}
	if parsed["type"] != "object" {
		t.Errorf("expected schema type 'object', got %v", parsed["type"])
	}
	props, ok := parsed["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties in schema")
	}
	if _, ok := props["id"]; !ok {
		t.Error("expected 'id' property in schema")
	}
	if _, ok := props["name"]; !ok {
		t.Error("expected 'name' property in schema")
	}
}

func TestGenerateSpec_EmptyResponseSchema(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/empty", ResponseSchema: ""},
	}

	data, err := GenerateSpec(endpoints, "Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc spec
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	op := doc.Paths["/empty"]["get"]
	resp := op.Responses["200"]
	mt := resp.Content[contentTypeJSON]

	var parsed map[string]any
	if err := json.Unmarshal(mt.Schema, &parsed); err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}
	if parsed["type"] != "object" {
		t.Errorf("expected default schema type 'object', got %v", parsed["type"])
	}
}

func TestGenerateSpec_InvalidJSONResponseSchema(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/bad", ResponseSchema: "not-json"},
	}

	data, err := GenerateSpec(endpoints, "Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc spec
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	op := doc.Paths["/bad"]["get"]
	resp := op.Responses["200"]
	mt := resp.Content[contentTypeJSON]

	var parsed map[string]any
	if err := json.Unmarshal(mt.Schema, &parsed); err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}
	if parsed["type"] != "string" {
		t.Errorf("expected wrapped schema type 'string', got %v", parsed["type"])
	}
	if parsed["example"] != "not-json" {
		t.Errorf("expected example 'not-json', got %v", parsed["example"])
	}
}

func TestGenerateSpec_EmptyEndpoints(t *testing.T) {
	data, err := GenerateSpec(nil, "Empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc spec
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if doc.OpenAPI != openAPIVersion {
		t.Errorf("expected openapi %q, got %q", openAPIVersion, doc.OpenAPI)
	}
	if doc.Info.Title != "Empty" {
		t.Errorf("expected title %q, got %q", "Empty", doc.Info.Title)
	}
	if len(doc.Paths) != 0 {
		t.Errorf("expected 0 paths, got %d", len(doc.Paths))
	}
}

func TestGenerateSpec_AllHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	endpoints := make([]model.Endpoint, len(methods))
	for i, m := range methods {
		endpoints[i] = model.Endpoint{Method: m, Path: "/resource", ResponseSchema: `{"type":"object"}`}
	}

	data, err := GenerateSpec(endpoints, "All Methods")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc spec
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	pi := doc.Paths["/resource"]
	expectedMethods := []string{"get", "post", "put", "delete", "patch"}
	for _, m := range expectedMethods {
		if _, ok := pi[m]; !ok {
			t.Errorf("expected %s operation on /resource", m)
		}
	}
}

func TestGenerateSpec_UnsupportedMethod(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "TRACE", Path: "/x", ResponseSchema: ""},
	}

	_, err := GenerateSpec(endpoints, "Test")
	if err == nil {
		t.Fatal("expected error for unsupported method TRACE")
	}
}

func TestGenerateSpec_OutputIsValidJSON(t *testing.T) {
	endpoints := []model.Endpoint{
		{Method: "GET", Path: "/a", ResponseSchema: `{"type":"string"}`},
		{Method: "POST", Path: "/b", ResponseSchema: ""},
	}

	data, err := GenerateSpec(endpoints, "JSON Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !json.Valid(data) {
		t.Error("output is not valid JSON")
	}
}

func TestGenerateSpec_RequiredOpenAPIFields(t *testing.T) {
	data, err := GenerateSpec(nil, "Fields Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	for _, field := range []string{"openapi", "info", "paths"} {
		if _, ok := raw[field]; !ok {
			t.Errorf("missing required OpenAPI field %q", field)
		}
	}
}
