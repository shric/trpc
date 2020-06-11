package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// VerifyOptions is all the command line options for the verify command.
type VerifyOptions struct {
	filter.Options `group:"filters"`
	ForceAll       bool `long:"force-all" description:"Really verify all torrents"`
}

func Verify(c *Command) {
	opts, ok := c.CommandOptions.(VerifyOptions)
	optionsCheck(ok)

	if len(c.PositionalArgs) == 0 && !opts.ForceAll {
		fmt.Fprintln(os.Stderr, "Use --force-all if you really want to verify all torrents!")
		return
	}

	ProcessTorrents(c.Client, opts.Options, c.PositionalArgs, []string{"name", "id"},
		func(torrent *transmissionrpc.Torrent) {
			if !c.CommonOptions.DryRun {
				err := c.Client.TorrentVerifyIDs([]int64{*torrent.ID})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
			}
			c.status("Verifying torrent", torrent)
		})
}
