package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// OpenAPI version constant.
const openAPIVersion = "3.0.0"

// spec is the top-level OpenAPI 3.0 document structure.
type spec struct {
	OpenAPI string              `json:"openapi"`
	Info    info                `json:"info"`
	Paths   map[string]pathItem `json:"paths"`
}

// info describes the API.
type info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// pathItem maps HTTP methods to operations for a given path.
type pathItem map[string]operation

// operation describes a single API operation.
type operation struct {
	Summary   string              `json:"summary"`
	Responses map[string]response `json:"responses"`
}

// response describes a single response from an API operation.
type response struct {
	Description string               `json:"description"`
	Content     map[string]mediaType `json:"content,omitempty"`
}

// mediaType describes the media type of a response.
type mediaType struct {
	Schema json.RawMessage `json:"schema"`
}

// defaultResponseSchema is used when the endpoint has no response schema defined.
var defaultResponseSchema = json.RawMessage(`{"type":"object"}`)

// contentTypeJSON is the standard JSON content type.
const contentTypeJSON = "application/json"

// GenerateSpec converts a slice of endpoint definitions into a valid OpenAPI 3.0.0
// JSON document. The title parameter is used for the spec's info.title field.
func GenerateSpec(endpoints []model.Endpoint, title string) ([]byte, error) {
	paths := make(map[string]pathItem, len(endpoints))

	for _, ep := range endpoints {
		method := strings.ToLower(ep.Method)
		if !isValidMethod(method) {
			return nil, fmt.Errorf("unsupported HTTP method: %q", ep.Method)
		}

		schema := parseResponseSchema(ep.ResponseSchema)

		op := operation{
			Summary: fmt.Sprintf("%s %s", strings.ToUpper(method), ep.Path),
			Responses: map[string]response{
				"200": {
					Description: "Successful response",
					Content: map[string]mediaType{
						contentTypeJSON: {Schema: schema},
					},
				},
			},
		}

		pi, exists := paths[ep.Path]
		if !exists {
			pi = make(pathItem)
			paths[ep.Path] = pi
		}
		pi[method] = op
	}

	doc := spec{
		OpenAPI: openAPIVersion,
		Info: info{
			Title:   title,
			Version: "1.0.0",
		},
		Paths: paths,
	}

	return json.MarshalIndent(doc, "", "  ")
}

// validMethods is the set of HTTP methods supported by OpenAPI path items.
var validMethods = map[string]bool{
	"get":    true,
	"post":   true,
	"put":    true,
	"delete": true,
	"patch":  true,
}

// isValidMethod checks whether the given lowercase method is a supported HTTP method.
func isValidMethod(method string) bool {
	return validMethods[method]
}

// parseResponseSchema converts a user-provided response schema string into a
// JSON Schema value suitable for embedding in the OpenAPI spec.
func parseResponseSchema(raw string) json.RawMessage {
	trimmed := strings.TrimSpace(raw)

	// Empty string → default object schema.
	if trimmed == "" {
		return defaultResponseSchema
	}

	// Valid JSON object → use directly.
	var obj map[string]any
	if err := json.Unmarshal([]byte(trimmed), &obj); err == nil {
		return json.RawMessage(trimmed)
	}

	// Invalid JSON → wrap as string example.
	wrapped, _ := json.Marshal(map[string]any{
		"type":    "string",
		"example": trimmed,
	})
	return json.RawMessage(wrapped)
}
