package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/util"
)

type rmOptions struct {
	torrentOptions
	filter.Options `group:"filters"`
	ForceAll       bool `long:"force-all" description:"Really allow all torrents to be removed"`
	Nuke           bool `long:"nuke" description:"Delete the data associated with the torrent"`
}

// Rm implements the rm command.
func Rm(c *Command) {
	opts, ok := c.Options.(rmOptions)
	optionsCheck(ok)

	if len(opts.Pos.Torrents) == 0 && !opts.ForceAll {
		fmt.Fprintln(os.Stderr, "Use --force-all if you really want to delete all torrents!")
		return
	}

	util.ProcessTorrents(c.Client, opts.Options, opts.Pos.Torrents, commonArgs[:],
		func(torrent *transmissionrpc.Torrent) {
			if !c.CommonOptions.DryRun {
				err := c.Client.TorrentRemove(&transmissionrpc.TorrentRemovePayload{
					IDs:             []int64{*torrent.ID},
					DeleteLocalData: opts.Nuke,
				})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
					return
				}
			}
			c.status("Removed torrent", torrent)
		}, nil, false)
}
