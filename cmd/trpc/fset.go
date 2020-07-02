package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"

	"github.com/shric/trpc/internal/util"
)

type fsetOptions struct {
	fileOptions
	Get      bool   `long:"get" short:"g" description:"Mark file for downloading"`
	NoGet    bool   `long:"noget" short:"G" description:"Mark file for not downloading"`
	Priority string `long:"priority" short:"p" description:"Set file priority" choice:"low" choice:"normal" choice:"high"`
}

func torrentFset(c *Command, files map[int64][]int64) {
	opts, ok := c.Options.(fsetOptions)
	optionsCheck(ok)

	for ID, fileIDs := range files {
		IDs := make([]int64, 1)
		IDs[0] = ID
		payload := &transmissionrpc.TorrentSetPayload{IDs: IDs}

		for _, fileID := range fileIDs {
			if opts.Get {
				payload.FilesWanted = append(payload.FilesWanted, fileID)
			}

			if opts.NoGet {
				payload.FilesUnwanted = append(payload.FilesUnwanted, fileID)
			}

			switch opts.Priority {
			case "high":
				payload.PriorityHigh = append(payload.PriorityHigh, fileID)
			case "normal":
				payload.PriorityNormal = append(payload.PriorityNormal, fileID)
			case "low":
				payload.PriorityLow = append(payload.PriorityLow, fileID)
			}
		}

		if !c.CommonOptions.DryRun {
			err := c.Client.TorrentSet(payload)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
	}
}

// Fset implements the which command (find which torrent a file/path is associated with.
func Fset(c *Command) {
	opts, ok := c.Options.(fsetOptions)
	optionsCheck(ok)

	finder := util.NewFinder(c.Client)

	// files maps torrent ID
	files := map[int64][]int64{}

	for _, f := range opts.Pos.Files {
		torrent, fileID := finder.Find(f)
		if torrent != nil {
			files[*torrent.ID] = append(files[*torrent.ID], fileID)
		} else {
			fmt.Fprintln(os.Stderr, "Couldn't find a torrent for", f)
		}
	}

	torrentFset(c, files)
}
