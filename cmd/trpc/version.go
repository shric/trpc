package cmd

import "fmt"

var (
	sha1ver   string
	buildTime string
	version   string
)

// Version prints the version number and build info.
func Version() {
	fmt.Printf("trpc version %s (%s) built at %s\n", version, sha1ver, buildTime)
}
