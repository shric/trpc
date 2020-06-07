package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/hekmon/transmissionrpc"
)

// AddOptions is all the command line options for the add command
type AddOptions struct {
	Paused bool `short:"p" long:"paused" description:"add torrent paused"`
}

// Add adds a new torrent by URL or file
func Add(client *transmissionrpc.Client, opts AddOptions, args []string) {
	for _, arg := range args {
		var torrent *transmissionrpc.Torrent

		url, err := url.Parse(arg)

		// Assume it's a file.
		if err != nil || url.Scheme == "" {
			fmt.Println("Treating as file.")
			filepath := arg
			var b64 string
			b64, err = transmissionrpc.File2Base64(filepath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "can't encode '%s' content as base64: %v", filepath, err)
			} else {
				torrent, err = client.TorrentAdd(&transmissionrpc.TorrentAddPayload{MetaInfo: &b64, Paused: &opts.Paused})
			}
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
