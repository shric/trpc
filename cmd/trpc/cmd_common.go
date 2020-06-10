package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

func getCanonicalFnames(fnames []string) (canonicalFnames map[string]int64) {
	canonicalFnames = make(map[string]int64)
	for _, fn := range fnames {
		canonicalPath, err := filepath.Abs(fn)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		canonicalFnames[canonicalPath] = -1
	}
	return
}

func getIncompleteDir(client *transmissionrpc.Client) (incompleteDir string, enabled bool) {
	session, err := client.SessionArgumentsGet()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if !*session.IncompleteDirEnabled {
		return
	}
	incompleteDir = *session.IncompleteDir
	return
}

// getids attempts to convert a list of torrent filenames to their corresponding ID
// numbers in transmission.
func getids(client *transmissionrpc.Client, fnames []string) (ids []int64) {
	// Let's do no work if given an empty list as this function is expensive
	if len(fnames) == 0 {
		return
	}
	canonicalFnames := getCanonicalFnames(fnames)
	paths := make([]string, 1, 2)
	incompleteDir, enabled := getIncompleteDir(client)
	if enabled {
		paths = append(paths, incompleteDir)
	}

	torrents, err := client.TorrentGet([]string{"id", "downloadDir", "name"}, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, torrent := range torrents {
		paths[0] = *torrent.DownloadDir
		for _, path := range paths {
			fullpath, err := filepath.Abs(filepath.Join(path, *torrent.Name))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			if canonicalFnames[fullpath] != 0 {
				canonicalFnames[fullpath] = *torrent.ID
				ids = append(ids, *torrent.ID)
			}
		}
	}
	for k, v := range canonicalFnames {
		if v == -1 {
			fmt.Fprintln(os.Stderr, "Did not find any torrent ID for", k)
		}
	}
	return
}

// ProcessTorrents runs the supplied function over all torrents matching the args and filters.
func ProcessTorrents(client *transmissionrpc.Client, filterOptions filter.Options, args []string, fields []string, do func(torrent *transmissionrpc.Torrent)) {
	ids := make([]int64, 0, len(args))

	f := filter.New(filterOptions)

	fields = append(fields, f.Args...)

	fnames := make([]string, 0, len(args))
	for _, strID := range args {
		if id, err := strconv.ParseInt(strID, 10, 64); err == nil {
			ids = append(ids, id)
		} else {
			fnames = append(fnames, strID)
		}
	}
	ids = append(ids, getids(client, fnames)...)
	torrents, err := client.TorrentGet(fields, ids)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, transmissionrpcTorrent := range torrents {
		if !f.CheckFilter(transmissionrpcTorrent) {
			continue
		}
		do(transmissionrpcTorrent)
	}
}
