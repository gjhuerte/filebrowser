package fbhttp

// securityHeaders are sent with every response, including the app page.
//
//   - nosniff stops browsers from guessing a different content type, e.g.
//     running an uploaded file as a script.
//   - SAMEORIGIN stops other sites from framing File Browser to trick users
//     into clicking (clickjacking). Same-origin framing, used for PDF previews,
//     keeps working.
//   - same-origin keeps file paths in URLs from leaking to other sites through
//     the Referer header.
var securityHeaders = map[string]string{
	"X-Content-Type-Options": "nosniff",
	"X-Frame-Options":        "SAMEORIGIN",
	"Referrer-Policy":        "same-origin",
}

// globalHeaders are also sent with every response. API answers and the app
// page must always be fetched fresh, since they depend on the signed-in user.
var globalHeaders = map[string]string{
	"Cache-Control": "no-cache, no-store, must-revalidate",
}
