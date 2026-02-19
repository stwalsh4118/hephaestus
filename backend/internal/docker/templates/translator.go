package templates

import (
	"fmt"

	"github.com/stwalsh4118/hephaestus/backend/internal/docker"
	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// portsRequired returns the number of host ports a service type needs.
func portsRequired(serviceType string) int {
	if serviceType == model.ServiceTypeRabbitMQ {
		return 2 // AMQP + management UI
	}
	return 1
}

// Translator converts a model.Diagram into an ordered slice of docker.ContainerConfig.
// Translator is not safe for concurrent use. Create separate instances for
// concurrent translations.
type Translator struct {
	registry  TemplateRegistry
	allocator *PortAllocator
}

// NewTranslator creates a Translator with the default registry and port allocator.
func NewTranslator() *Translator {
	return &Translator{
		registry:  NewRegistry(),
		allocator: NewPortAllocator(DefaultMinPort, DefaultMaxPort),
	}
}

// Translate converts a diagram into an ordered slice of container configs.
// The order respects dependency ordering (infrastructure before application).
// The port allocator is reset for each translation call.
func (t *Translator) Translate(diagram model.Diagram) ([]docker.ContainerConfig, error) {
	if len(diagram.Nodes) == 0 {
		return nil, nil
	}

	// Validate all node types.
	for _, node := range diagram.Nodes {
		if _, ok := t.registry[node.Type]; !ok {
			return nil, fmt.Errorf("unsupported service type %q for node %q", node.Type, node.ID)
		}
	}

	// Resolve startup order.
	order, err := ResolveDependencies(diagram.Nodes, diagram.Edges)
	if err != nil {
		return nil, fmt.Errorf("resolve dependencies: %w", err)
	}

	// Build node lookup.
	nodeMap := make(map[string]model.DiagramNode, len(diagram.Nodes))
	for _, n := range diagram.Nodes {
		nodeMap[n.ID] = n
	}

	// Reset port allocator for fresh deployment.
	t.allocator.Reset()

	// Build configs in dependency order.
	configs := make([]docker.ContainerConfig, 0, len(order))
	for _, nodeID := range order {
		node := nodeMap[nodeID]
		tmpl := t.registry[node.Type]

		n := portsRequired(node.Type)
		ports, err := t.allocator.AllocateN(n)
		if err != nil {
			return nil, fmt.Errorf("allocate ports for node %q: %w", nodeID, err)
		}

		cfg, err := tmpl.Build(node, ports[0], ports[1:]...)
		if err != nil {
			return nil, fmt.Errorf("build config for node %q: %w", nodeID, err)
		}

		configs = append(configs, cfg)
	}

	return configs, nil
}
