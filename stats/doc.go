// Package stats computes the Home dashboard's numbers (file and folder
// counts, sizes by type, recent and largest files) on the server, as a
// background job.
//
// Walking a large tree takes time, so the HTTP layer starts a scan with
// [Service.Start] and clients poll [Service.Status] until it finishes. Only
// the small [Summary] travels to the browser, never the full file list.
package stats
