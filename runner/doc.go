// Package runner executes the commands configured for file events (the
// "before_upload", "after_delete", … hooks) and parses command lines.
//
// Commands run either directly or through the configured shell. Event details
// are passed as environment variables (FILE, SCOPE, TRIGGER, USERNAME,
// DESTINATION) so that file names never become part of a shell script.
package runner
