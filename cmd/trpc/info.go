package cmd

import (
	"fmt"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/util"
)

type infoOptions struct {
	torrentOptions
	filter.Options `group:"filters"`
}

// Info provides a list of files for all or selected torrents.
func Info(c *Command) {
	opts, ok := c.Options.(infoOptions)
	optionsCheck(ok)
	util.ProcessTorrents(c.Client, opts.Options, opts.Pos.Torrents, append(commonArgs[:], "files", "priorities", "wanted", "hashString", "magnetLink"),
		func(transmissionrpcTorrent *transmissionrpc.Torrent) {
			fmt.Println(info(transmissionrpcTorrent))
		}, nil, false)
}

func info(t *transmissionrpc.Torrent) string {
	return infoGeneral(t) + "\n" + infoTransfer(t)
}

func infoGeneral(t *transmissionrpc.Torrent) string {
	return "NAME\n" +
		fmt.Sprintf("  Id: %d\n", *t.ID) +
		fmt.Sprintf("  Name: %s\n", *t.Name) +
		fmt.Sprintf("  Hash: %s\n", *t.HashString) +
		fmt.Sprintf("  Magnet: %s\n", *t.MagnetLink)
}

func infoTransfer(t *transmissionrpc.Torrent) string {
	return "TRANSFER\n" +
		fmt.Sprintf("  State: %v\n", *t.Status)
}
