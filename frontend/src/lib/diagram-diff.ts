import type { DiagramJson } from "@/lib/topology";
import type { CanvasNode } from "@/types/canvas";

/**
 * Returns true if the current canvas nodes differ from the last deployed diagram.
 * Compares node IDs using set equality — detects additions and removals.
 */
export function hasDiagramChanges(
  nodes: CanvasNode[],
  lastDeployed: DiagramJson | null
): boolean {
  if (!lastDeployed) return false;
  const currentIds = new Set(nodes.map((n) => n.id));
  const deployedIds = new Set(lastDeployed.nodes.map((n) => n.id));
  // Size mismatch catches pure removals; iteration catches swaps/additions.
  if (currentIds.size !== deployedIds.size) return true;
  for (const id of currentIds) {
    if (!deployedIds.has(id)) return true;
  }
  return false;
}
