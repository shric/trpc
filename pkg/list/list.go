// Package list provides the list command and its formatting.
package list

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/shric/go-trpc/pkg/torrent"

	"github.com/hekmon/transmissionrpc"
)

func format(pointerTorrent *transmissionrpc.Torrent) (string) {
	result := torrent.NewFrom(pointerTorrent)
	format := `{{printf "%4v" .ID}} {{.Error}} {{printf "%5.1f" .Percent}}% {{printf "%11s" .Size}} {{printf "%8s" .Eta}} {{.Name}}`
	var tpl bytes.Buffer
	tmpl := template.Must(template.New("list").Parse(format))
	tmpl.Execute(&tpl, result)
	return tpl.String()
}

// List provides a list of all or selected torrents.
func List(client *transmissionrpc.Client) {


	torrents, err := client.TorrentGet([]string{
		"name", "recheckProgress", "sizeWhenDone", "rateUpload", "eta", "id",
		"leftUntilDone", "recheckProgress", "error", "rateDownload",
		"status", "trackers", "bandwidthPriority", "uploadedEver",
		"downloadDir", "addedDate", "doneDate", "startDate",
		"isFinished",
	}, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		for _, torrent := range torrents {
			formattedTorrent := format(torrent)
			fmt.Println(formattedTorrent)
		}
	}

}
