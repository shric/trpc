package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// StartOptions is all the command line options for the stop command.
type StartOptions struct {
	filter.Options `group:"filters"`
	Now            bool `long:"now" description:"Start torrent now, bypassing the queue"`
}

// Start starts torrents.
func Start(c *Command) {
	opts, ok := c.Options.(StartOptions)
	optionsCheck(ok)

	startFunc := c.Client.TorrentStartIDs

	if opts.Now {
		startFunc = c.Client.TorrentStartNowIDs
	}

	ProcessTorrents(c.Client, opts.Options, c.PositionalArgs, []string{"name", "id"},
		func(torrent *transmissionrpc.Torrent) {
			if !c.CommonOptions.DryRun {
				err := startFunc([]int64{*torrent.ID})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
			c.status("Started torrent", torrent)
		})
}

// StopOptions is all the command line options for the stop command.
type StopOptions struct {
	filter.Options `group:"filters"`
}

// Stop stops torrents.
func Stop(c *Command) {
	opts, ok := c.Options.(StopOptions)
	optionsCheck(ok)
	ProcessTorrents(c.Client, opts.Options, c.PositionalArgs, []string{"name", "id"},
		func(torrent *transmissionrpc.Torrent) {
			if !c.CommonOptions.DryRun {
				err := c.Client.TorrentStopIDs([]int64{*torrent.ID})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
			c.status("Stopped torrent", torrent)
		})
}
