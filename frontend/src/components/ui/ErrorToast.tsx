"use client";

import { useEffect } from "react";

const AUTO_DISMISS_MS = 5000;

interface ErrorToastProps {
  message: string;
  onDismiss: () => void;
}

export function ErrorToast({ message, onDismiss }: ErrorToastProps) {
  useEffect(() => {
    const timer = setTimeout(onDismiss, AUTO_DISMISS_MS);
    return () => clearTimeout(timer);
  }, [message, onDismiss]);

  return (
    <div className="fixed bottom-4 right-4 z-50 flex max-w-sm items-start gap-2 rounded-lg border border-red-300 bg-red-50 px-4 py-3 shadow-lg">
      <p className="flex-1 text-sm text-red-800">{message}</p>
      <button
        type="button"
        onClick={onDismiss}
        className="ml-2 text-red-400 hover:text-red-600"
        aria-label="Dismiss error"
      >
        &times;
      </button>
    </div>
  );
}
