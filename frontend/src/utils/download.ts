/*
 * Download a file into a blob URL while reporting progress, so previews
 * can show "12% · 3.4 MB of 28 MB" instead of a blank screen.
 */

export interface DownloadProgress {
  loaded: number;
  /** 0 when the server didn't send a size and no fallback was given. */
  total: number;
}

export async function downloadWithProgress(
  url: string,
  opts: {
    signal?: AbortSignal;
    /** Used when the response has no Content-Length (e.g. thumbnails). */
    fallbackTotal?: number;
    onProgress?: (p: DownloadProgress) => void;
  } = {}
): Promise<string> {
  const res = await fetch(url, {
    signal: opts.signal,
    credentials: "same-origin",
  });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);

  const header = Number(res.headers.get("Content-Length")) || 0;
  // Compressed responses report the compressed length; ignore it then
  const encoded = !!res.headers.get("Content-Encoding");
  const total = (!encoded && header) || opts.fallbackTotal || 0;
  const type = res.headers.get("Content-Type") || "";

  if (!res.body) {
    const blob = await res.blob();
    opts.onProgress?.({ loaded: blob.size, total: blob.size });
    return URL.createObjectURL(blob);
  }

  const reader = res.body.getReader();
  const chunks: BlobPart[] = [];
  let loaded = 0;
  opts.onProgress?.({ loaded, total });

  for (;;) {
    const { done, value } = await reader.read();
    if (done) break;
    chunks.push(value as Uint8Array<ArrayBuffer>);
    loaded += value.length;
    opts.onProgress?.({ loaded, total: Math.max(total, loaded) });
  }

  return URL.createObjectURL(new Blob(chunks, { type }));
}

/** Fraction (0–1) of the media that is buffered, or null if unknown. */
export function bufferedFraction(el: HTMLMediaElement): number | null {
  const d = el.duration;
  if (!d || !isFinite(d) || el.buffered.length === 0) return null;
  // Use the range around the playhead; fall back to the furthest one
  const t = el.currentTime;
  let end = 0;
  for (let i = 0; i < el.buffered.length; i++) {
    if (el.buffered.start(i) <= t + 0.5 && el.buffered.end(i) >= t) {
      end = el.buffered.end(i);
      break;
    }
    end = Math.max(end, el.buffered.end(i));
  }
  return Math.min(1, end / d);
}
