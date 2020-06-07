// Package list provides the list command and its formatting.
package list

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/shric/go-trpc/pkg/config"
	"github.com/shric/go-trpc/pkg/torrent"

	"github.com/hekmon/transmissionrpc"
)

// Options declares all the command line arguments for list
// These should probably be moved to their own package because they're generic filter options.
type Options struct {
	Incomplete bool `short:"i" long:"incomplete" description:"only incomplete torrents"`
}

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

type filter struct {
	predicate func(*torrent.Torrent, string) bool
	set       interface{}
}

func checkFilter(result *torrent.Torrent, opts Options) bool {
	var Filters = []filter{
		{
			predicate: func(t *torrent.Torrent, v string) bool {
				return t.LeftUntilDone > 0
			},
			set: &opts.Incomplete,
		},
	}
	match := true

	for _, f := range Filters {
		switch v := f.set.(type) {
		case *bool:
			if *v && !f.predicate(result, "") {
				match = false
			}
		}
	}
	return match
}

// List provides a list of all or selected torrents
func List(client *transmissionrpc.Client, opts Options, args []string) {

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
	} else {
		for _, transmissionrpcTorrent := range torrents {
			result := torrent.NewFrom(transmissionrpcTorrent, conf)
			if !checkFilter(result, opts) {
				continue
			}
			total.Size += result.Size
			total.Up += result.Up
			total.Down += result.Down
			formattedTorrent := format(result, conf)
			fmt.Println(formattedTorrent)
		}
		formattedTotal := format(&total, conf)
		fmt.Println(formattedTotal)
	}

}
