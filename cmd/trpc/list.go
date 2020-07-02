package cmd

import (
	"fmt"
	"os"

	"github.com/shric/trpc/internal/config"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/torrent"
	"github.com/shric/trpc/internal/util"

	"github.com/hekmon/transmissionrpc"
	"github.com/slongfield/pyfmt"
)

func format(torrent *torrent.Torrent, _ *config.Config) string {
	format := "{ID:4}{Error:1} {Pct:5}%  {Size:6.1f} {SizeSuffix:<3} {Eta:<8} {Up:>7} {Down:>7}" +
		" {Ratio:>6.1f}  {Priority:6} {Trackershortname:4}   {Name}"
	return pyfmt.Must(format, torrent)
}

type listOptions struct {
	torrentOptions
	filter.Options `group:"filters"`
	NoTotals       bool   `short:"n" long:"no-totals" description:"suppress output of totals"`
	Sort           string `long:"sort" description:"sort" choice:"size" choice:"name" choice:"id" choice:"ratio" choice:"have" choice:"progress" choice:"uploaded" choice:"age"`
	Reverse        bool   `short:"r" long:"reverse" description:"reverse sort order"`
}

// List provides a list of all or selected torrents.
func List(c *Command) {
	opts, ok := c.Options.(listOptions)
	optionsCheck(ok)

	total := torrent.NewForTotal()

	total.Error = " "

	conf := config.ReadConfig()

	if c.CommonOptions.DryRun {
		fmt.Fprintln(os.Stderr, "--dry-run has no effect on list as list doesn't change state")
	}

	linePrinted := false

	var sortField *string
	if opts.Sort != "" {
		sortField = &opts.Sort
	}

	util.ProcessTorrents(c.Client, opts.Options, opts.Pos.Torrents, commonArgs[:],
		func(transmissionrpcTorrent *transmissionrpc.Torrent) {
			result := torrent.NewFrom(transmissionrpcTorrent, conf)
			total.UpdateTotal(result)

			formattedTorrent := format(result, conf)
			fmt.Println(formattedTorrent)
			linePrinted = true
		}, sortField, opts.Reverse)

	if !opts.NoTotals && linePrinted {
		formattedTotal := format(total, conf)
		fmt.Println(formattedTotal)
	}
}
