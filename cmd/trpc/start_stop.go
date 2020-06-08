package cmd

import (
	"fmt"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// StartOptions is all the command line options for the stop command
type StartOptions struct {
	Now            bool `long:"now" short:"n" description:"Start torrent now, bypassing the queue"`
	filter.Options `group:"filters"`
}

// Start starts torrents.
func Start(client *transmissionrpc.Client, opts StartOptions, args []string) {
	startFunc := client.TorrentStartIDs

	if opts.Now {
		startFunc = client.TorrentStartNowIDs
	}

	ProcessTorrents(client, opts.Options, args, []string{"name", "id"}, func(torrent *transmissionrpc.Torrent) {
		startFunc([]int64{*torrent.ID})
		fmt.Printf("Started torrent %d: %s\n", *torrent.ID, *torrent.Name)
	})
	return
}

// StopOptions is all the command line options for the stop command
type StopOptions struct {
	filter.Options `group:"filters"`
}

// Stop stops torrents.
func Stop(client *transmissionrpc.Client, opts StopOptions, args []string) {
	ProcessTorrents(client, opts.Options, args, []string{"name", "id"}, func(torrent *transmissionrpc.Torrent) {
		client.TorrentStopIDs([]int64{*torrent.ID})
		fmt.Printf("Stopped torrent %d: %s\n", *torrent.ID, *torrent.Name)
	})
	return
}
