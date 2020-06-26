package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/shric/trpc/internal/fileutils"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/util"
)

type renameOptions struct {
	ID         int64 `long:"torrent-id" short:"t" description:"Use this torrent ID instead of inferring from local filesystem"`
	Positional struct {
		Oldname string `positional-arg-name:"old" description:"the old name or ID of the torrent's path"`
		Newname string `positional-arg-name:"new" description:"the new name of the torrent's path"`
	} `positional-args:"true"`
}

// Rename renames a torrent path or file.
func Rename(c *Command) {
	opts, ok := c.Options.(renameOptions)
	optionsCheck(ok)

	finder := util.NewFinder(c.Client)

	oldname := fileutils.RealPath(opts.Positional.Oldname)
	newname := fileutils.RealPath(opts.Positional.Newname)

	var torrent *transmissionrpc.Torrent

	if opts.ID != 0 {
		torrents, err := c.Client.TorrentGet([]string{"id", "downloadDir", "name"}, []int64{opts.ID})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Torrent ID %d not found.", opts.ID)
			os.Exit(1)
		}

		torrent = torrents[0]
	} else {
		torrent, _ = finder.Find(oldname)
	}

	if torrent == nil {
		fmt.Fprintln(os.Stderr, "Couldn't determine associated torrent from ", oldname)
		os.Exit(1)
	}

	realDownloadDir := fileutils.RealPath(*torrent.DownloadDir) + "/"

	oldname = strings.Replace(oldname, realDownloadDir, "", 1)
	newname = path.Base(newname)

	var err error

	if !c.CommonOptions.DryRun {
		err = c.Client.TorrentRenamePath(*torrent.ID, oldname, newname)
	}

	if err != nil {
		fmt.Println("Rename: err: ", err)
		os.Exit(1)
	} else {
		c.statusf("Renamed %s to %s", oldname, newname)
	}
}
