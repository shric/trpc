package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hekmon/cunits/v2"

	"github.com/shric/trpc/internal/torrent"

	"github.com/slongfield/pyfmt"

	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/util"

	"github.com/hekmon/transmissionrpc"
)

func fileInfo(t *transmissionrpc.Torrent) string {
	var s string

	result := torrent.NewFrom(t, nil)

	s = pyfmt.Must("{ID}: {Name}:\n", result)
	formatString := "%3s: %4s %-8s %3s %10s  %s"
	s += fmt.Sprintf(formatString+"\n",
		"#", "Done", "Priority", "Get", "Size", "Name")

	for i, f := range t.Files {
		var get string
		if t.Wanted[i] {
			get = "Yes"
		} else {
			get = "No "
		}

		s += fmt.Sprintf(formatString+"\n",
			strconv.Itoa(i),
			strconv.FormatInt(int64(100.0*torrent.FileProgress(t, int64(i))), 10)+"%",
			torrent.FilePriority(t, int64(i)),
			get,
			cunits.ImportInByte(float64(f.Length)).GetHumanSizeRepresentation(),
			strings.Replace(f.Name, *t.Name+"/", "", 1))
	}

	return s
}

type filesOptions struct {
	torrentOptions
	filter.Options `group:"filters"`
}

// Files provides a list of files for all or selected torrents.
func Files(c *Command) {
	opts, ok := c.Options.(filesOptions)
	optionsCheck(ok)
	util.ProcessTorrents(c.Client, opts.Options, opts.Pos.Torrents, append(commonArgs[:], "files", "priorities", "wanted"),
		func(transmissionrpcTorrent *transmissionrpc.Torrent) {
			fmt.Println(fileInfo(transmissionrpcTorrent))
		}, nil, false)
}
