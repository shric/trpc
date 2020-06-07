// Package list provides the list command and its formatting.
package list

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/shric/go-trpc/pkg/config"
	"github.com/shric/go-trpc/pkg/filter"
	"github.com/shric/go-trpc/pkg/torrent"

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

// Options defines all the command line options for list.
type Options struct {
	filter.Options `group:"filters"`
	NoTotals       bool `short:"n" long:"no-totals" description:"suppress output of totals"`
}

// List provides a list of all or selected torrents
func List(client *transmissionrpc.Client, opts Options, args []string) {
	f := filter.New(opts.Options)

	torrents, err := client.TorrentGet([]string{
		"name", "recheckProgress", "sizeWhenDone", "rateUpload", "eta", "id",
		"leftUntilDone", "recheckProgress", "error", "rateDownload",
		"status", "trackers", "bandwidthPriority", "uploadedEver",
		"downloadDir", "addedDate", "doneDate", "startDate",
		"isFinished",
	}, nil)

	var total torrent.Torrent
	total.Error = " "
	conf := config.ReadConfig()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, transmissionrpcTorrent := range torrents {
		if !f.CheckFilter(transmissionrpcTorrent) {
			continue
		}
		result := torrent.NewFrom(transmissionrpcTorrent, conf)
		total.Size += result.Size
		total.Up += result.Up
		total.Down += result.Down
		formattedTorrent := format(result, conf)
		fmt.Println(formattedTorrent)
	}
	formattedTotal := format(&total, conf)
	fmt.Println(formattedTotal)
}
