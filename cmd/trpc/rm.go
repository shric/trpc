package cmd

import (
	"fmt"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// RmOptions is all the command line options for the rm command
type RmOptions struct {
	filter.Options `group:"filters"`
	ForceAll       bool `long:"force-all" description:"Really allow all torrents to be removed"`
	Nuke           bool `long:"nuke" description:"Delete the data associated with the torrent"`
}

// Rm implements the rm command.
func Rm(client *transmissionrpc.Client, opts RmOptions, args []string) {
	if len(args) == 0 && !opts.ForceAll {
		return
	}
	ProcessTorrents(client, opts.Options, args, []string{"name", "id"}, func(torrent *transmissionrpc.Torrent) {
		client.TorrentRemove(&transmissionrpc.TorrentRemovePayload{
			IDs:             []int64{*torrent.ID},
			DeleteLocalData: opts.Nuke,
		})
		fmt.Printf("Removed torrent %d: %s\n", *torrent.ID, *torrent.Name)
	})
}
