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

func format(torrent *torrent.Torrent, conf *config.Config) string {
	format := `{{printf "%4v" .ID}} {{.Error}} {{printf "%6.1f" .Percent}}%  {{printf "%12s" .Size}} {{printf "%-8s" .Eta}} {{printf "%7.1f" .Up}} {{printf "%7.1f" .Down}} {{printf "%6.1f" .Ratio}} {{printf "%-6s" .Priority}}  {{printf "%-4s" .Trackershortname}}  {{.Name}}`
	var tpl bytes.Buffer
	tmpl := template.Must(template.New("list").Parse(format))
	err := tmpl.Execute(&tpl, torrent)
	if err != nil {
		os.Exit(1)
	}
	return tpl.String()
}

// ListOptions defines all the command line options for list.
type ListOptions struct {
	filter.Options `group:"filters"`
	NoTotals       bool `short:"n" long:"no-totals" description:"suppress output of totals"`
}

// List provides a list of all or selected torrents
func List(client *transmissionrpc.Client, opts ListOptions, args []string) {
	var total torrent.Torrent
	total.Error = " "
	conf := config.ReadConfig()

	ProcessTorrents(client, opts.Options, args, []string{
		"name", "recheckProgress", "sizeWhenDone", "rateUpload", "eta", "id",
		"leftUntilDone", "recheckProgress", "error", "rateDownload",
		"status", "trackers", "bandwidthPriority", "uploadedEver",
		"downloadDir", "addedDate", "doneDate", "startDate",
		"isFinished",
	}, func(transmissionrpcTorrent *transmissionrpc.Torrent) {
		result := torrent.NewFrom(transmissionrpcTorrent, conf)
		total.Size += result.Size
		total.Up += result.Up
		total.Down += result.Down
		formattedTorrent := format(result, conf)
		fmt.Println(formattedTorrent)
	})
	formattedTotal := format(&total, conf)
	fmt.Println(formattedTotal)
}
