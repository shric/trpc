package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/util"
)

// VerifyOptions is all the command line options for the verify command.
type verifyOptions struct {
	torrentOptions
	filter.Options `group:"filters"`
	ForceAll       bool `long:"force-all" description:"Really verify all torrents"`
}

// Verify is a command that verifies (hash checks) the selected torrents.
func Verify(c *Command) {
	opts, ok := c.Options.(verifyOptions)
	optionsCheck(ok)

	if len(opts.Pos.Torrents) == 0 && !opts.ForceAll {
		fmt.Fprintln(os.Stderr, "Use --force-all if you really want to verify all torrents!")
		return
	}

	util.ProcessTorrents(c.Client, opts.Options, opts.Pos.Torrents, commonArgs[:],
		func(torrent *transmissionrpc.Torrent) {
			if !c.CommonOptions.DryRun {
				err := c.Client.TorrentVerifyIDs([]int64{*torrent.ID})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
			}
			c.status("Verifying torrent", torrent)
		}, nil, false)
}
