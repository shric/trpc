package add

import (
	"fmt"
	"net/url"

	"github.com/hekmon/transmissionrpc"
)

// Options is all the command line options for the add command
type Options struct {
	Paused bool `short:"p" long:"paused" description:"add torrent paused"`
}

// Add adds a new torrent by URL or file
func Add(client *transmissionrpc.Client, opts Options, args []string) {
	for _, arg := range args {
		var torrent *transmissionrpc.Torrent

		url, err := url.Parse(arg)

		// Assume it's a file.
		if err != nil || url.Scheme == "" {
			fmt.Println("Treating as file.")
			// TODO: provide paused files.
			// https://github.com/hekmon/transmissionrpc/issues/11
			torrent, err = client.TorrentAddFile(arg)
		} else { // It's a URL, pass it to transmission.
			payload := &transmissionrpc.TorrentAddPayload{
				Filename: &arg,
				Paused:   &opts.Paused,
			}
			torrent, err = client.TorrentAdd(payload)

		}
		if err != nil {
			fmt.Println("Add: err: ", err)
		} else {
			fmt.Printf("Added torrent with ID %d: %s\n", *torrent.ID, *torrent.Name)
		}

	}
}
