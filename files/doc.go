// Package files reads file metadata for the API: [FileInfo] describes a file
// or folder, its type (detected from the extension or, optionally, the
// content), its listing and checksums.
//
// It also provides [ScopedFs], the filesystem every request goes through. It
// confines a user to their scope and refuses to follow symbolic links whose
// target lies outside it.
package files
