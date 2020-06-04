package add

import (
	"fmt"
	"net/url"

	"github.com/hekmon/transmissionrpc"
)

// Add adds a new torrent by URL or file
func Add(client *transmissionrpc.Client, args []string) {
	for _, arg := range args {
		var torrent *transmissionrpc.Torrent

		url, err := url.Parse(arg)

		// Assume it's a file.
		if err != nil || url.Scheme == "" {
			fmt.Println("Treating as file.")
			torrent, err = client.TorrentAddFile(arg)
		} else { // It's a URL, pass it to transmission.
			payload := &transmissionrpc.TorrentAddPayload{
				Filename: &arg,
			}
			torrent, err = client.TorrentAdd(payload)

		}
		if err != nil {
			fmt.Println("err: ", err)
		}

		fmt.Printf("%+v\n", torrent)
		fmt.Println(arg)
	}
}
