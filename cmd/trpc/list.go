package cmd

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/shric/trpc/internal/config"
	"github.com/shric/trpc/internal/filter"
	"github.com/shric/trpc/internal/torrent"

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
	filter.Options `group:"filters"`
	NoTotals       bool `short:"n" long:"no-totals" description:"suppress output of totals"`
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

	torrent.ProcessTorrents(c.Client, opts.Options, c.PositionalArgs, []string{
		"name", "recheckProgress", "sizeWhenDone", "rateUpload", "eta", "id",
		"leftUntilDone", "recheckProgress", "error", "rateDownload",
		"status", "trackers", "bandwidthPriority", "uploadedEver",
		"downloadDir", "addedDate", "doneDate", "startDate",
		"isFinished",
	}, func(transmissionrpcTorrent *transmissionrpc.Torrent) {
		result := torrent.NewFrom(transmissionrpcTorrent, conf)
		total.UpdateTotal(result)

		formattedTorrent := format(result, conf)
		fmt.Println(formattedTorrent)
	})

	formattedTotal := format(total, conf)
	fmt.Println(formattedTotal)
}
