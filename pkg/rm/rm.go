package rm

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/go-trpc/pkg/filter"
)

// Options is all the command line options for the rm command
type Options struct {
	filter.Options
	ForceAll bool `long:"force-all" description:"Really allow all torrents to be removed"`
	Nuke     bool `long:"nuke" description:"Delete the data associated with the torrent"`
}

// Rm implements the rm command.
func Rm(client *transmissionrpc.Client, opts Options, args []string) {
	if len(args) == 0 && !opts.ForceAll {
		return
	}
	ids := make([]int64, 0, len(args))

	f := filter.New(opts.Options)

	fields := []string{"name", "id"}
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
		client.TorrentRemove(&transmissionrpc.TorrentRemovePayload{
			IDs:             []int64{*transmissionrpcTorrent.ID},
			DeleteLocalData: opts.Nuke,
		})
		fmt.Printf("Removed torrent %d: %s\n", *transmissionrpcTorrent.ID, *transmissionrpcTorrent.Name)
	}
}
