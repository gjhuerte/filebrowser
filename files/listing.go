package files

import (
	"sort"
	"strings"

	"github.com/maruel/natural"
)

// Listing is a collection of files.
type Listing struct {
	Items    []*FileInfo `json:"items"`
	NumDirs  int         `json:"numDirs"`
	NumFiles int         `json:"numFiles"`
	Sorting  Sorting     `json:"sorting"`
}

// ApplySort sorts the items in place by l.Sorting.By ("name", "size" or
// "modified") in the direction given by l.Sorting.Asc.
func (l Listing) ApplySort() {
	// Check '.Order' to know how to sort
	if !l.Sorting.Asc {
		switch l.Sorting.By {
		case "name":
			sort.Sort(sort.Reverse(byName(l)))
		case "size":
			sort.Sort(sort.Reverse(bySize(l)))
		case "modified":
			sort.Sort(sort.Reverse(byModified(l)))
		default:
			// If not one of the above, do nothing
			return
		}
	} else { // If we had more Orderings we could add them here
		switch l.Sorting.By {
		case "name":
			sort.Sort(byName(l))
		case "size":
			sort.Sort(bySize(l))
		case "modified":
			sort.Sort(byModified(l))
		default:
			sort.Sort(byName(l))
			return
		}
	}
}

// byName, bySize and byModified implement sort.Interface for a Listing. All
// three keep folders together before (or, reversed, after) files.
type byName Listing
type bySize Listing
type byModified Listing

// Len implements sort.Interface.
func (l byName) Len() int {
	return len(l.Items)
}

// Swap implements sort.Interface.
func (l byName) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}

// Less implements sort.Interface, comparing names naturally ("file2" before
// "file10") and ignoring case.
func (l byName) Less(i, j int) bool {
	if l.Items[i].IsDir && !l.Items[j].IsDir {
		return l.Sorting.Asc
	}

	if !l.Items[i].IsDir && l.Items[j].IsDir {
		return !l.Sorting.Asc
	}

	return natural.Less(strings.ToLower(l.Items[j].Name), strings.ToLower(l.Items[i].Name))
}

// Len implements sort.Interface.
func (l bySize) Len() int {
	return len(l.Items)
}

// Swap implements sort.Interface.
func (l bySize) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}

// directoryOffset is added to folder sizes so folders always sort before files.
const directoryOffset = -1 << 31 // = math.MinInt32

// Less implements sort.Interface, ordering by size.
func (l bySize) Less(i, j int) bool {
	iSize, jSize := l.Items[i].Size, l.Items[j].Size
	if l.Items[i].IsDir {
		iSize = directoryOffset + iSize
	}
	if l.Items[j].IsDir {
		jSize = directoryOffset + jSize
	}
	return iSize < jSize
}

// Len implements sort.Interface.
func (l byModified) Len() int {
	return len(l.Items)
}

// Swap implements sort.Interface.
func (l byModified) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}

// Less implements sort.Interface, ordering by modification time.
func (l byModified) Less(i, j int) bool {
	iModified, jModified := l.Items[i].ModTime, l.Items[j].ModTime
	return iModified.Sub(jModified) < 0
}
