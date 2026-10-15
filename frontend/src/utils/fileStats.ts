/*
 * Summarise a recursive directory listing for the Home dashboard.
 * Runs in a single pass so it stays fast on large trees.
 */

export interface RecursiveEntryLike {
  path: string;
  name: string;
  size: number;
  modified: string;
  isDir: boolean;
}

export type FileTypeKey =
  | "images"
  | "videos"
  | "audio"
  | "documents"
  | "archives"
  | "other";

export interface ScanSummary {
  files: number;
  folders: number;
  bytes: number;
  byType: { key: FileTypeKey; bytes: number; count: number }[];
  recent: RecursiveEntryLike[];
  largest: RecursiveEntryLike[];
}

export const TYPE_ICONS: Record<FileTypeKey, string> = {
  images: "image",
  videos: "movie",
  audio: "audiotrack",
  documents: "description",
  archives: "inventory_2",
  other: "insert_drive_file",
};

const EXTENSIONS: Record<Exclude<FileTypeKey, "other">, string[]> = {
  images: [
    "jpg",
    "jpeg",
    "png",
    "gif",
    "webp",
    "svg",
    "bmp",
    "tif",
    "tiff",
    "heic",
    "heif",
    "avif",
    "ico",
    "raw",
    "cr2",
    "nef",
    "arw",
    "dng",
  ],
  videos: [
    "mp4",
    "mkv",
    "mov",
    "avi",
    "webm",
    "m4v",
    "wmv",
    "flv",
    "mpg",
    "mpeg",
    "3gp",
  ],
  audio: [
    "mp3",
    "flac",
    "wav",
    "aac",
    "ogg",
    "oga",
    "m4a",
    "opus",
    "wma",
    "aiff",
    "alac",
  ],
  documents: [
    "pdf",
    "doc",
    "docx",
    "xls",
    "xlsx",
    "ppt",
    "pptx",
    "odt",
    "ods",
    "odp",
    "txt",
    "md",
    "csv",
    "rtf",
    "epub",
    "json",
    "xml",
    "html",
    "htm",
  ],
  archives: [
    "zip",
    "tar",
    "gz",
    "tgz",
    "bz2",
    "xz",
    "7z",
    "rar",
    "zst",
    "iso",
    "dmg",
  ],
};

const BY_EXTENSION = new Map<string, FileTypeKey>();
for (const [key, exts] of Object.entries(EXTENSIONS)) {
  for (const ext of exts) BY_EXTENSION.set(ext, key as FileTypeKey);
}

export function fileType(name: string): FileTypeKey {
  const dot = name.lastIndexOf(".");
  if (dot <= 0) return "other";
  return BY_EXTENSION.get(name.slice(dot + 1).toLowerCase()) ?? "other";
}

/** Keep the top `n` items of `list` by `score`, highest first. */
function pushTop<T>(list: T[], item: T, n: number, score: (x: T) => number) {
  const s = score(item);
  if (list.length === n && s <= score(list[n - 1])) return;
  let i = list.length;
  while (i > 0 && score(list[i - 1]) < s) i--;
  list.splice(i, 0, item);
  if (list.length > n) list.pop();
}

export function summarize(entries: RecursiveEntryLike[], top = 6): ScanSummary {
  const totals = new Map<FileTypeKey, { bytes: number; count: number }>();
  const recent: RecursiveEntryLike[] = [];
  const largest: RecursiveEntryLike[] = [];
  let files = 0;
  let folders = 0;
  let bytes = 0;

  for (const e of entries) {
    if (e.isDir) {
      folders++;
      continue;
    }
    files++;
    bytes += e.size;
    const key = fileType(e.name);
    const t = totals.get(key) ?? { bytes: 0, count: 0 };
    t.bytes += e.size;
    t.count++;
    totals.set(key, t);
    pushTop(recent, e, top, (x) => Date.parse(x.modified) || 0);
    pushTop(largest, e, top, (x) => x.size);
  }

  const byType = [...totals.entries()]
    .map(([key, v]) => ({ key, ...v }))
    .sort((a, b) =>
      a.key === "other" ? 1 : b.key === "other" ? -1 : b.bytes - a.bytes
    );

  return { files, folders, bytes, byType, recent, largest };
}
