package cmd

import (
	"fmt"

	"github.com/hekmon/transmissionrpc"
)

var (
	sha1ver   string
	buildTime string
	version   string
)

// Version prints the version number and build info.
func Version(c *Command) {
	fmt.Printf("trpc version %s (%s) built at %s\n", version, sha1ver, buildTime)

	ok, serverVersion, serverMinimumVersion, err := c.Client.RPCVersion()
	if err != nil {
		panic(err)
	}

	if !ok {
		panic(fmt.Sprintf(
			"Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): min v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
	}

	fmt.Printf("Client library is built against RPC version v%d\n", transmissionrpc.RPCVersion)
	fmt.Printf("Remote transmission RPC version v%d\n", serverVersion)
}
