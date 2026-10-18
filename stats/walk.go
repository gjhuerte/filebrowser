package stats

import (
	"context"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

// Options configures a scan.
type Options struct {
	// Root is the real directory on disk that maps to the virtual path "/".
	Root string
	// Allow reports whether a virtual path may be counted. Returning false
	// for a folder skips everything inside it. Nil allows everything.
	Allow func(virtual string) bool
	// Workers is how many folders are read at the same time. Zero picks a
	// default suited to both local disks and network shares.
	Workers int
}

// Progress is updated live while a scan runs.
type Progress struct {
	Files   atomic.Int64
	Folders atomic.Int64
}

// Scan walks everything under opts.Root and returns its Summary.
//
// It is built for large trees:
//   - folders are listed with os.ReadDir, which reports each entry's type
//     without an extra system call; only regular files are stat'ed, for
//     their size and modification time
//   - several folders are read concurrently, which matters most on network
//     storage and SSDs
//   - symbolic links are never followed, so a scan cannot leave Root or loop
//   - only totals and the top entries are kept, so memory use stays flat no
//     matter how many files there are
//
// Folders that cannot be read are skipped. Scan stops early, returning
// ctx.Err(), when ctx is cancelled.
func Scan(ctx context.Context, opts Options, progress *Progress) (*Summary, error) {
	if progress == nil {
		progress = &Progress{}
	}
	workers := opts.Workers
	if workers <= 0 {
		workers = max(4, min(16, runtime.NumCPU()*2))
	}
	allow := opts.Allow
	if allow == nil {
		allow = func(string) bool { return true }
	}

	w := &walker{
		ctx:      ctx,
		allow:    allow,
		progress: progress,
		slots:    make(chan struct{}, workers-1),
		total:    newTally(),
	}
	w.dir(opts.Root, "/")
	w.wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return w.total.summary(), nil
}

type walker struct {
	ctx      context.Context
	allow    func(string) bool
	progress *Progress
	// slots limits extra goroutines; when none is free, the folder is read
	// on the current goroutine instead, so the walk never blocks.
	slots chan struct{}
	wg    sync.WaitGroup

	mu    sync.Mutex
	total *tally
}

// dir counts the folder at real (virtual path virtual) and walks its
// subfolders, in parallel where a slot is free.
func (w *walker) dir(real, virtual string) {
	if w.ctx.Err() != nil {
		return
	}

	entries, err := os.ReadDir(real)
	if err != nil {
		return // unreadable folders are skipped
	}

	local := newTally()
	var subdirs []string
	for _, e := range entries {
		name := e.Name()
		v := path.Join(virtual, name)
		if !w.allow(v) {
			continue
		}

		mode := e.Type()
		switch {
		case mode&fs.ModeSymlink != 0:
			// Not followed: the target may be outside the scope.
			continue
		case mode.IsDir():
			local.folders++
			subdirs = append(subdirs, name)
		case mode.IsRegular():
			info, err := e.Info()
			if err != nil {
				continue // removed while scanning
			}
			local.addFile(v, info.Size(), info.ModTime())
		default:
			// Devices, sockets and pipes have no meaningful size.
		}
	}

	w.progress.Files.Add(local.files)
	w.progress.Folders.Add(local.folders)
	w.mu.Lock()
	w.total.merge(local)
	w.mu.Unlock()

	for _, name := range subdirs {
		subReal := filepath.Join(real, name)
		subVirtual := path.Join(virtual, name)
		select {
		case w.slots <- struct{}{}:
			w.wg.Add(1)
			go func() {
				defer func() {
					<-w.slots
					w.wg.Done()
				}()
				w.dir(subReal, subVirtual)
			}()
		default:
			w.dir(subReal, subVirtual)
		}
	}
}
