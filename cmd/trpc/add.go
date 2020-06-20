package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/config"
	"github.com/shric/trpc/internal/utils"
)

type addOptions struct {
	Paused      bool   `short:"p" long:"paused" description:"add torrent paused"`
	DownloadDir string `short:"d" long:"download-dir" dscription:"download directory"`
}

// Add adds a new torrent by URL or file.
func Add(c *Command) {
	opts, ok := c.Options.(addOptions)
	optionsCheck(ok)

	// These silence G601: Implicit memory aliasing in for loop. (gosec)
	var argCopy string

	var dummyID int64

	if len(c.PositionalArgs) == 0 {
		fmt.Fprintln(os.Stderr, "Please supply at least one file or URL")
		os.Exit(1)
	}

	conf := config.ReadConfig()
	if opts.DownloadDir == "" && conf.Settings.Has("default_download_dir") {
		opts.DownloadDir = conf.Settings.Get("default_download_dir").(string)
	}

	for _, arg := range c.PositionalArgs {
		var torrent *transmissionrpc.Torrent

		url, err := url.Parse(arg)

		payload := transmissionrpc.TorrentAddPayload{
			Paused: &opts.Paused,
		}

		if opts.DownloadDir != "" {
			realDownloadDir := utils.RealPath(opts.DownloadDir)
			payload.DownloadDir = &realDownloadDir
		}

		// Assume it's a file.
		if err != nil || url.Scheme == "" {
			fmt.Println("Treating as file.")

			b64, err := transmissionrpc.File2Base64(arg)
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
		} else {
			// Fill it with something for dry-run
			dummyID = 0
			argCopy = arg
			torrent = &transmissionrpc.Torrent{
				ID:   &dummyID,
				Name: &argCopy}
		}

		if err != nil {
			fmt.Println("Add: err: ", err)
			os.Exit(1)
		} else {
			c.status("Added torrent with ID", torrent)
		}
	}
}
