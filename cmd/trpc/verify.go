package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// VerifyOptions is all the command line options for the verify command
type VerifyOptions struct {
	filter.Options `group:"filters"`
	ForceAll       bool `long:"force-all" description:"Really verify all torrents"`
}

// Verify implements the verify command (hash check torrents)
func Verify(client *transmissionrpc.Client, opts VerifyOptions, args []string) {
	if len(args) == 0 && !opts.ForceAll {
		fmt.Fprintln(os.Stderr, "Use --force-all if you really want to verify all torrents!")
		return
	}
	ProcessTorrents(client, opts.Options, args, []string{"name", "id"}, func(torrent *transmissionrpc.Torrent) {
		err := client.TorrentVerifyIDs([]int64{*torrent.ID})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Printf("Verifying torrent %d: %s\n", *torrent.ID, *torrent.Name)
	})
}
