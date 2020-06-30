package cmd

import (
	"fmt"
	"os"

	"github.com/shric/trpc/internal/fileutils"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/util"
)

type moveOptions struct {
	Positional struct {
		Files []string `positional-arg-name:"torrent" description:"0 or more torrents with destination at the end"`
	} `positional-args:"true"`
	filter.Options `group:"filters"`
	ForceAll       bool `long:"force-all" description:"Really move all torrents"`
}

func getFnamesAndDest(args []string) (fnames []string, dest string) {
	dest = fileutils.RealPath(args[len(args)-1])
	fnames = args[:len(args)-1]

	return
}

// Move implements the verify command (hash check torrents).
func Move(c *Command) {
	opts, ok := c.Options.(moveOptions)
	optionsCheck(ok)

	if len(opts.Positional.Files) == 0 {
		fmt.Fprintln(os.Stderr, "move: Destination required")
		return
	}

	if len(opts.Positional.Files) == 1 && !opts.ForceAll {
		fmt.Fprintln(os.Stderr, "Use --force-all if you really want to move all torrents")
		return
	}

	fnames, destination := getFnamesAndDest(opts.Positional.Files)
	util.ProcessTorrents(c.Client, opts.Options, fnames, commonArgs[:], func(torrent *transmissionrpc.Torrent) {
		if !c.CommonOptions.DryRun {
			err := c.Client.TorrentSetLocation(*torrent.ID, destination, true)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		c.status("Moving torrent", torrent)
	}, nil, false)
}
