package stats

import (
	"path"
	"sort"
	"strings"
	"time"
)

// TopN is how many entries the recent and largest lists keep.
const TopN = 6

// TypeKey groups files into the categories shown on the Home dashboard.
type TypeKey string

// The categories a file can fall into, decided by its extension.
const (
	TypeImages    TypeKey = "images"
	TypeVideos    TypeKey = "videos"
	TypeAudio     TypeKey = "audio"
	TypeDocuments TypeKey = "documents"
	TypeArchives  TypeKey = "archives"
	TypeOther     TypeKey = "other"
)

// extensions maps a lower-case extension (without the dot) to its category.
// Keep in sync with frontend/src/utils/fileStats.ts.
var extensions = func() map[string]TypeKey {
	groups := map[TypeKey][]string{
		TypeImages: {"jpg", "jpeg", "png", "gif", "webp", "svg", "bmp", "tif", "tiff",
			"heic", "heif", "avif", "ico", "raw", "cr2", "nef", "arw", "dng"},
		TypeVideos: {"mp4", "mkv", "mov", "avi", "webm", "m4v", "wmv", "flv", "mpg",
			"mpeg", "3gp"},
		TypeAudio: {"mp3", "flac", "wav", "aac", "ogg", "oga", "m4a", "opus", "wma",
			"aiff", "alac"},
		TypeDocuments: {"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "odt", "ods",
			"odp", "txt", "md", "csv", "rtf", "epub", "json", "xml", "html", "htm"},
		TypeArchives: {"zip", "tar", "gz", "tgz", "bz2", "xz", "7z", "rar", "zst", "iso",
			"dmg"},
	}
	m := map[string]TypeKey{}
	for key, exts := range groups {
		for _, ext := range exts {
			m[ext] = key
		}
	}
	return m
}()

// FileType returns the category for a file name.
func FileType(name string) TypeKey {
	dot := strings.LastIndexByte(name, '.')
	if dot <= 0 {
		return TypeOther
	}
	if key, ok := extensions[strings.ToLower(name[dot+1:])]; ok {
		return key
	}
	return TypeOther
}

// Entry is one file in the recent or largest list. Path is the virtual path
// inside the user's scope, "/"-separated.
type Entry struct {
	Path     string    `json:"path"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
	IsDir    bool      `json:"isDir"`
}

// TypeTotal is the size and number of files in one category.
type TypeTotal struct {
	Key   TypeKey `json:"key"`
	Bytes int64   `json:"bytes"`
	Count int64   `json:"count"`
}

// Summary describes everything under a folder. Its JSON shape matches the
// frontend's ScanSummary.
type Summary struct {
	Files   int64       `json:"files"`
	Folders int64       `json:"folders"`
	Bytes   int64       `json:"bytes"`
	ByType  []TypeTotal `json:"byType"`
	Recent  []Entry     `json:"recent"`
	Largest []Entry     `json:"largest"`
}

// tally accumulates a Summary. It is not safe for concurrent use: each
// worker keeps its own and they are merged at the end.
type tally struct {
	files, folders, bytes int64
	types                 map[TypeKey]*TypeTotal
	recent, largest       []Entry
}

func newTally() *tally {
	return &tally{types: map[TypeKey]*TypeTotal{}}
}

func (t *tally) addFile(virtual string, size int64, mod time.Time) {
	t.files++
	t.bytes += size
	name := path.Base(virtual)
	key := FileType(name)
	tt := t.types[key]
	if tt == nil {
		tt = &TypeTotal{Key: key}
		t.types[key] = tt
	}
	tt.Bytes += size
	tt.Count++

	e := Entry{Path: virtual, Name: name, Size: size, Modified: mod}
	t.recent = pushTop(t.recent, e, func(a, b Entry) bool { return a.Modified.After(b.Modified) })
	t.largest = pushTop(t.largest, e, func(a, b Entry) bool { return a.Size > b.Size })
}

// merge adds o into t.
func (t *tally) merge(o *tally) {
	t.files += o.files
	t.folders += o.folders
	t.bytes += o.bytes
	for key, ot := range o.types {
		tt := t.types[key]
		if tt == nil {
			tt = &TypeTotal{Key: key}
			t.types[key] = tt
		}
		tt.Bytes += ot.Bytes
		tt.Count += ot.Count
	}
	for _, e := range o.recent {
		t.recent = pushTop(t.recent, e, func(a, b Entry) bool { return a.Modified.After(b.Modified) })
	}
	for _, e := range o.largest {
		t.largest = pushTop(t.largest, e, func(a, b Entry) bool { return a.Size > b.Size })
	}
}

// summary returns the finished Summary. Categories are ordered by size, with
// "other" always last.
func (t *tally) summary() *Summary {
	s := &Summary{
		Files:   t.files,
		Folders: t.folders,
		Bytes:   t.bytes,
		ByType:  make([]TypeTotal, 0, len(t.types)),
		Recent:  append([]Entry{}, t.recent...),
		Largest: append([]Entry{}, t.largest...),
	}
	for _, tt := range t.types {
		s.ByType = append(s.ByType, *tt)
	}
	sort.Slice(s.ByType, func(i, j int) bool {
		a, b := s.ByType[i], s.ByType[j]
		if (a.Key == TypeOther) != (b.Key == TypeOther) {
			return b.Key == TypeOther
		}
		if a.Bytes != b.Bytes {
			return a.Bytes > b.Bytes
		}
		return a.Key < b.Key
	})
	return s
}

// pushTop inserts e into list, kept sorted by before and capped at TopN.
func pushTop(list []Entry, e Entry, before func(a, b Entry) bool) []Entry {
	if len(list) == TopN && !before(e, list[TopN-1]) {
		return list
	}
	i := sort.Search(len(list), func(i int) bool { return before(e, list[i]) })
	list = append(list, Entry{})
	copy(list[i+1:], list[i:])
	list[i] = e
	if len(list) > TopN {
		list = list[:TopN]
	}
	return list
}
