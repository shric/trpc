package cmd

import (
	"fmt"
	"os"

	"github.com/shric/trpc/internal/util"
)

type whichOptions struct {
	fileOptions
	Missing bool `long:"missing" description:"Show only unassociated files/paths"`
}

// Which implements the which command (find which torrent a file/path is associated with.
func Which(c *Command) {
	opts, ok := c.Options.(whichOptions)
	optionsCheck(ok)

	finder := util.NewFinder(c.Client)

	for _, f := range opts.Positional.Files {
		torrent, fileID := finder.Find(f)
		if torrent != nil {
			fmt.Printf("%s belongs to torrent %d: %s (File ID %d)\n",
				f, *torrent.ID, *torrent.Name, fileID)
		} else {
			fmt.Fprintln(os.Stderr, "Couldn't find a torrent for", f)
		}
	}
}
