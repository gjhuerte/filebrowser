// Package cmd holds the filebrowser command-line interface: the root command
// that starts the server, and the subcommands that manage configuration,
// users, rules and command hooks in the database.
//
// Most subcommands open the Bolt database exclusively, so they cannot run
// while the server is using the same database file.
package cmd
