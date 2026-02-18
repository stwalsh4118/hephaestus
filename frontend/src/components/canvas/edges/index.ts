import type { EdgeTypes } from "@xyflow/react";

import { EDGE_TYPE_LABELED } from "@/constants/canvas";

import { LabeledEdge } from "./LabeledEdge";

export const edgeTypes: EdgeTypes = {
  [EDGE_TYPE_LABELED]: LabeledEdge,
};
