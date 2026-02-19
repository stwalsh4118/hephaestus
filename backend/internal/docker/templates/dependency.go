package templates

import (
	"errors"
	"fmt"
	"sort"

	"github.com/stwalsh4118/hephaestus/backend/internal/model"
)

// Service type priority levels — lower number means start first.
const (
	PriorityInfrastructure = 0 // postgresql, redis, rabbitmq
	PriorityApplication    = 1 // nginx, api-service
)

// ErrCyclicDependency is returned when the dependency graph contains a cycle.
var ErrCyclicDependency = errors.New("cyclic dependency detected")

// servicePriority returns the startup priority for a service type.
func servicePriority(serviceType string) int {
	switch serviceType {
	case model.ServiceTypePostgreSQL, model.ServiceTypeRedis, model.ServiceTypeRabbitMQ:
		return PriorityInfrastructure
	default:
		return PriorityApplication
	}
}

// ResolveDependencies performs a topological sort on diagram nodes using
// Kahn's algorithm. Edges define dependencies: if edge source→target exists,
// then target should start before source (target is a dependency of source).
// Nodes with equal in-degree are ordered by service type priority.
func ResolveDependencies(nodes []model.DiagramNode, edges []model.DiagramEdge) ([]string, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	// Build node index for type lookups.
	nodeTypes := make(map[string]string, len(nodes))
	for _, n := range nodes {
		nodeTypes[n.ID] = n.Type
	}

	// Build adjacency list and in-degree map.
	// edge source→target means source depends on target,
	// so target → source in the adjacency list (target must start first).
	adj := make(map[string][]string, len(nodes))
	inDegree := make(map[string]int, len(nodes))
	for _, n := range nodes {
		adj[n.ID] = nil
		inDegree[n.ID] = 0
	}

	for _, e := range edges {
		// Skip edges referencing unknown nodes.
		if _, ok := nodeTypes[e.Source]; !ok {
			continue
		}
		if _, ok := nodeTypes[e.Target]; !ok {
			continue
		}
		adj[e.Target] = append(adj[e.Target], e.Source)
		inDegree[e.Source]++
	}

	// Seed queue with zero in-degree nodes, sorted by priority.
	queue := make([]string, 0)
	for _, n := range nodes {
		if inDegree[n.ID] == 0 {
			queue = append(queue, n.ID)
		}
	}
	sortByPriority(queue, nodeTypes)

	// Kahn's algorithm.
	result := make([]string, 0, len(nodes))
	for len(queue) > 0 {
		// Pop first element.
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Collect neighbors whose in-degree drops to zero.
		var newReady []string
		for _, neighbor := range adj[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				newReady = append(newReady, neighbor)
			}
		}
		if len(newReady) > 0 {
			sortByPriority(newReady, nodeTypes)
			queue = append(queue, newReady...)
		}
	}

	if len(result) != len(nodes) {
		return nil, fmt.Errorf("%w: processed %d of %d nodes", ErrCyclicDependency, len(result), len(nodes))
	}

	return result, nil
}

// sortByPriority sorts node IDs by service type priority (infrastructure first),
// with stable alphabetical ordering for ties within the same priority.
func sortByPriority(ids []string, nodeTypes map[string]string) {
	sort.SliceStable(ids, func(i, j int) bool {
		pi := servicePriority(nodeTypes[ids[i]])
		pj := servicePriority(nodeTypes[ids[j]])
		if pi != pj {
			return pi < pj
		}
		return ids[i] < ids[j]
	})
}
