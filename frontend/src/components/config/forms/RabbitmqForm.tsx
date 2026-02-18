"use client";

import { useCanvasStore } from "@/store/canvas-store";
import { ST_RABBITMQ } from "@/types/canvas";
import type { RabbitMQConfig } from "@/types/canvas";

const DEFAULT_VHOST = "/";

interface RabbitmqFormProps {
  nodeId: string;
  config?: RabbitMQConfig;
}

export function RabbitmqForm({ nodeId, config }: RabbitmqFormProps) {
  const updateNodeConfig = useCanvasStore((state) => state.updateNodeConfig);

  const vhost = config?.vhost ?? DEFAULT_VHOST;

  return (
    <div className="space-y-4">
      <div>
        <label className="block text-xs font-medium text-slate-600">Virtual Host</label>
        <input
          type="text"
          value={vhost}
          onChange={(e) =>
            updateNodeConfig(nodeId, {
              type: ST_RABBITMQ,
              vhost: e.target.value,
            })
          }
          placeholder="e.g. /"
          className="mt-1 w-full rounded border border-slate-300 px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        />
      </div>
    </div>
  );
}
