// Package auth implements the ways a user can sign in to File Browser.
//
// Each method satisfies [Auther]:
//   - [JSONAuth] checks a username and password (with optional reCAPTCHA).
//   - [ProxyAuth] trusts a username header set by a reverse proxy, creating
//     the user on first sight.
//   - [HookAuth] asks an external command to approve, reject or defer.
//   - [NoAuth] signs everyone in as the first user.
//
// The method in use is stored by [Storage] and selected in the settings.
package auth
