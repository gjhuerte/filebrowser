import type { ScanSummary } from "@/utils/fileStats";
import { fetchJSON } from "./utils";

/** The server's state of the Home dashboard scan (see stats/service.go). */
export interface SummaryStatus {
  scanning: boolean;
  summary: ScanSummary | null;
  scannedAt?: string;
  durationMs?: number;
  progress: { files: number; folders: number };
  error?: string;
}

/** Read the last result and whether a scan is running. */
export function status(signal?: AbortSignal) {
  return fetchJSON<SummaryStatus>("/api/summary", { signal });
}

/** Start a scan (or join the one already running). Returns at once. */
export function start(signal?: AbortSignal) {
  return fetchJSON<SummaryStatus>("/api/summary", { method: "POST", signal });
}
