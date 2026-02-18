import type { NodeTypes } from "@xyflow/react";

import { ServiceNode } from "./ServiceNode";

export const NODE_TYPE_SERVICE = "service-node";

export const nodeTypes: NodeTypes = {
  [NODE_TYPE_SERVICE]: ServiceNode,
};
