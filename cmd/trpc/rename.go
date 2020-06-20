package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/shric/trpc/internal/torrent"
	"github.com/shric/trpc/internal/utils"
)

type renameOptions struct {
}

// Rename renames a torrent path or file.
func Rename(c *Command) {
	_, ok := c.Options.(renameOptions)
	optionsCheck(ok)

	finder := torrent.NewFinder(c.Client)

	if len(c.PositionalArgs) != 2 {
		fmt.Println("Rename requires oldname newname")
	}

	oldname := utils.RealPath(c.PositionalArgs[0])
	newname := utils.RealPath(c.PositionalArgs[1])

	torrent, _ := finder.Find(oldname)
	if torrent == nil {
		fmt.Fprintln(os.Stderr, "Couldn't determine associated torrent from ", oldname)
		os.Exit(1)
	}

	realDownloadDir := utils.RealPath(*torrent.DownloadDir) + "/"

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
		message := fmt.Sprintf("Renamed %s to %s", oldname, newname)
		c.status(message, torrent)
	}
}
