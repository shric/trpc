package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/util"
)

type errorsOptions struct {
	torrentOptions
	filter.Options `group:"filters"`
}

// Errors provides a list of all or selected torrents.
func Errors(c *Command) {
	opts, ok := c.Options.(errorsOptions)
	optionsCheck(ok)

	if c.CommonOptions.DryRun {
		fmt.Fprintln(os.Stderr, "--dry-run has no effect on errors as errors doesn't change state")
	}

	util.ProcessTorrents(c.Client, opts.Options, opts.Positional.Torrents, commonArgs[:],
		func(torrent *transmissionrpc.Torrent) {
			if *torrent.Error != 0 {
				fmt.Printf("ID: %5d %s:\n\t%s\n", *torrent.ID, *torrent.Name, *torrent.ErrorString)
			}
		}, nil, false)
}
