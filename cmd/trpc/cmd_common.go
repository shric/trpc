package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// ProcessTorrents runs the supplied function over all torrents matching the args and filters.
func ProcessTorrents(client *transmissionrpc.Client, filterOptions filter.Options, args []string, fields []string, do func(torrent *transmissionrpc.Torrent)) {
	ids := make([]int64, 0, len(args))

	f := filter.New(filterOptions)

	for _, arg := range f.Args {
		fields = append(fields, arg)
	}

	for _, strID := range args {
		if id, err := strconv.ParseInt(strID, 10, 64); err == nil {
			ids = append(ids, id)
		} else {
			fmt.Println("invalid id: ", strID)
		}
	}
	torrents, err := client.TorrentGet(fields, ids)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, transmissionrpcTorrent := range torrents {
		if !f.CheckFilter(transmissionrpcTorrent) {
			continue
		}
		do(transmissionrpcTorrent)
	}
	return
}
