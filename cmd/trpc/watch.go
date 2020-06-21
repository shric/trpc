package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/shric/trpc/internal/torrent"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/util"
)

const delayMillis = 200
const KiB = 1024

type watchOptions struct {
	filter.Options `group:"filters"`
}

func Watch(c *Command) {
	opts, ok := c.Options.(watchOptions)
	optionsCheck(ok)

	if c.CommonOptions.DryRun {
		fmt.Fprintln(os.Stderr, "--dry-run has no effect on watch as watch doesn't change state")
	}

	IDs := make([]int64, 0)

	util.ProcessTorrents(c.Client, opts.Options, c.PositionalArgs, commonArgs[:],
		func(torrent *transmissionrpc.Torrent) {
			if *torrent.Status == transmissionrpc.TorrentStatusStopped {
				return
			}
			if *torrent.LeftUntilDone == 0 {
				return
			}
			IDs = append(IDs, *torrent.ID)
		}, nil, false)

	if len(IDs) == 0 {
		return
	}

	var done bool

	for {
		torrents, err := c.Client.TorrentGet([]string{"leftUntilDone", "sizeWhenDone", "id", "name", "rateDownload", "recheckProgress"}, IDs)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Torrent get error: ", err)
			os.Exit(1)
		}

		done = true

		for _, t := range torrents {
			percent := torrent.Progress(t)

			var percentStr string

			var rateStr string

			if *t.LeftUntilDone != 0 {
				done = false
				percentStr = fmt.Sprintf("%5.2f%%", percent)
				rateStr = fmt.Sprintf("%9.2f KiB/s", float64(*t.RateDownload)/KiB)
			} else {
				percentStr = "Done"
				rateStr = ""
			}

			fmt.Printf("%5d: %6s %15s %s\n", *t.ID, percentStr, rateStr, *t.Name)
		}

		if done {
			break
		}

		time.Sleep(delayMillis * time.Millisecond)

		for range torrents {
			fmt.Print("\033[F")
		}
	}
}
