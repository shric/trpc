package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/hekmon/transmissionrpc"
)

// AddOptions is all the command line options for the add command.
type AddOptions struct {
	Paused bool `short:"p" long:"paused" description:"add torrent paused"`
}

// Add adds a new torrent by URL or file.
func Add(c *Command) {
	opts, ok := c.CommandOptions.(AddOptions)
	optionsCheck(ok)
	// This silences G601: Implicit memory aliasing in for loop. (gosec)
	var argCopy string
	for _, arg := range c.PositionalArgs {
		var torrent *transmissionrpc.Torrent

		url, err := url.Parse(arg)

		payload := transmissionrpc.TorrentAddPayload{
			Paused: &opts.Paused,
		}

		// Assume it's a file.
		if err != nil || url.Scheme == "" {
			fmt.Println("Treating as file.")
			var b64 string
			b64, err = transmissionrpc.File2Base64(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "can't encode '%s' content as base64: %v", arg, err)
				return
			}
			payload.MetaInfo = &b64
		} else { // It's a URL, pass it to transmission.
			argCopy = arg
			payload.Filename = &argCopy
		}
		if !c.CommonOptions.DryRun {
			torrent, err = c.Client.TorrentAdd(&payload)
		}
		if err != nil {
			fmt.Println("Add: err: ", err)
		} else {
			c.status("Added torrent with ID", torrent)
		}
	}
}
