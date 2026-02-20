import type { DiagramJson } from "@/lib/topology";
import type { DeployResponse, DeployStatusResponse } from "@/types/deploy";

const API_BASE =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

/**
 * Deploy a diagram by sending it to POST /api/deploy.
 */
export async function deployDiagram(
  diagram: DiagramJson
): Promise<DeployResponse> {
  const res = await fetch(`${API_BASE}/api/deploy`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(diagram),
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? `Deploy failed: ${res.status}`);
  }

  return res.json();
}

/**
 * Teardown all deployed containers via DELETE /api/deploy.
 */
export async function teardownDiagram(): Promise<DeployResponse> {
  const res = await fetch(`${API_BASE}/api/deploy`, {
    method: "DELETE",
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? `Teardown failed: ${res.status}`);
  }

  return res.json();
}

/**
 * Get the current deploy status via GET /api/deploy/status.
 */
export async function getDeployStatus(): Promise<DeployStatusResponse> {
  const res = await fetch(`${API_BASE}/api/deploy/status`);

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? `Status check failed: ${res.status}`);
  }

  return res.json();
}

/**
 * Update a deployed diagram via PUT /api/deploy.
 */
export async function updateDeploy(
  diagram: DiagramJson
): Promise<DeployStatusResponse> {
  const res = await fetch(`${API_BASE}/api/deploy`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(diagram),
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? `Update deploy failed: ${res.status}`);
  }

  return res.json();
}
