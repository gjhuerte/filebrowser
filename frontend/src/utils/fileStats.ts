/*
 * File categories for the Home dashboard and search.
 *
 * The Home numbers themselves are computed on the server (stats/summary.go),
 * which keeps its own copy of this extension list: keep the two in sync.
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
