// Package fbhttp is File Browser's HTTP layer: the router built by
// [NewHandler], authentication middleware, and the handlers behind the REST
// API, resumable uploads (tus), previews, public shares, search and the
// command terminal.
//
// Every handler has the [handleFunc] signature and is wrapped by handle, which
// loads the settings, applies the common response headers and turns the
// returned status and error into a response.
package fbhttp
