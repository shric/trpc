package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// RmOptions is all the command line options for the rm command.
type RmOptions struct {
	filter.Options `group:"filters"`
	ForceAll       bool `long:"force-all" description:"Really allow all torrents to be removed"`
	Nuke           bool `long:"nuke" description:"Delete the data associated with the torrent"`
}

// Rm implements the rm command.
func Rm(c *Command) {
	opts, ok := c.CommandOptions.(RmOptions)
	optionsCheck(ok)
	if len(c.PositionalArgs) == 0 && !opts.ForceAll {
		fmt.Fprintln(os.Stderr, "Use --force-all if you really want to delete all torrents!")
		return
	}
	ProcessTorrents(c.Client, opts.Options, c.PositionalArgs, []string{"name", "id"}, func(torrent *transmissionrpc.Torrent) {
		if !c.CommonOptions.DryRun {
			err := c.Client.TorrentRemove(&transmissionrpc.TorrentRemovePayload{
				IDs:             []int64{*torrent.ID},
				DeleteLocalData: opts.Nuke,
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		c.status("Removed torrent", torrent)
	})
}
