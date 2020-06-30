package cmd

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/shric/trpc/internal/config"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/torrent"
	"github.com/shric/trpc/internal/util"

	"github.com/hekmon/transmissionrpc"
)

func format(torrent *torrent.Torrent, _ *config.Config) string {
	format := `{{printf "%4s" .ID}} ` +
		`{{.Error}} ` +
		`{{printf "%6.1f" .Percent}}%  ` +
		`{{printf "%12s" .SizeWhenDone}} ` +
		`{{printf "%-8s" .Eta}} ` +
		`{{printf "%8s" .Up}} ` +
		`{{printf "%8s" .Down}} ` +
		`{{printf "%6.1f" .Ratio}} ` +
		`{{printf "%-6s" .Priority}}  ` +
		`{{printf "%-4s" .Trackershortname}}  ` +
		`{{.Name}}`

	var tpl bytes.Buffer

	tmpl := template.Must(template.New("list").Parse(format))

	err := tmpl.Execute(&tpl, torrent)
	if err != nil {
		panic(err)
	}

	return tpl.String()
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
