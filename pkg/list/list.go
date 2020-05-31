// Package list provides the list command and its formatting.
package list

import (
	"fmt"
	"os"
	"text/template"

	"github.com/hekmon/transmissionrpc"
)

// List provides a list of all or selected torrents.
func List(client *transmissionrpc.Client) {

	format := `{{printf "%-4v" (Derefint64 .ID)}}{{.Name}}`
	tmpl := template.Must(template.New("list").Funcs(template.FuncMap{
		"Derefint64": func(i *int64) int64 { return *i },
	}).Parse(format))

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
			tmpl.Execute(os.Stdout, *torrent)
			fmt.Println()
		}
	}

}
