package util

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/shric/trpc/internal/fileutils"

	"github.com/shric/trpc/internal/torrent"

	"github.com/shric/trpc/internal/config"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

func getAbsoluteFnames(fnames []string) (absoluteFnames map[string]int64) {
	absoluteFnames = make(map[string]int64)

	for _, fn := range fnames {
		realPath := fileutils.RealPath(fn)
		absoluteFnames[realPath] = -1
	}

	return
}

func getIncompleteDir(client *transmissionrpc.Client) *string {
	session, err := client.SessionArgumentsGet()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !*session.IncompleteDirEnabled {
		return nil
	}

	return session.IncompleteDir
}

// getids attempts to convert a list of torrent filenames to their corresponding ID
// numbers in transmission.
func getids(client *transmissionrpc.Client, fnames []string) []int64 {
	// Let's do no work if given an empty list as this function is expensive
	if len(fnames) == 0 {
		return nil
	}

	canonicalFnames := getAbsoluteFnames(fnames)
	paths := make([]string, 1, 2)
	incompleteDir := getIncompleteDir(client)

	if incompleteDir != nil {
		paths = append(paths, *incompleteDir)
	}

	torrents, err := client.TorrentGet([]string{"id", "downloadDir", "name"}, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	var ids []int64

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

	return ids
}

func sortTorrents(torrents []*transmissionrpc.Torrent, sortField string, reverse bool) {
	sort.SliceStable(torrents, func(i, j int) bool {
		x := torrents[i]
		y := torrents[j]
		if reverse {
			x, y = y, x
		}
		switch sortField {
		case "age":
			return torrent.Age(x) < torrent.Age(y)
		case "have":
			return torrent.Have(x) < torrent.Have(y)
		case "id":
			return *x.ID < *y.ID
		case "name":
			return strings.ToLower(*x.Name) < strings.ToLower(*y.Name)
		case "progress":
			return torrent.Progress(x) < torrent.Progress(y)
		case "ratio":
			return torrent.Ratio(x) < torrent.Ratio(y)
		case "size":
			return x.SizeWhenDone.Byte() < y.SizeWhenDone.Byte()
		case "uploaded":
			return *x.UploadedEver < *y.UploadedEver
		default:
			return *x.ID < *y.ID
		}
	})
}

// ProcessTorrents runs the supplied function over all torrents matching the args and filters.
func ProcessTorrents(client *transmissionrpc.Client, filterOptions filter.Options, args []string,
	fields []string, do func(torrent *transmissionrpc.Torrent), sortField *string, reverse bool,
) {
	ids := make([]int64, 0, len(args))

	conf := config.ReadConfig()
	f := filter.New(filterOptions, conf)

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
	// Something was specified as args but nothing could be converted to an ID.
	if len(ids) == 0 && len(args) > 0 {
		os.Exit(1)
	}

	torrents, err := client.TorrentGet(fields, ids)

	if sortField != nil {
		sortTorrents(torrents, *sortField, reverse)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, transmissionrpcTorrent := range torrents {
		if !f.CheckFilter(transmissionrpcTorrent) {
			continue
		}

		do(transmissionrpcTorrent)
	}
}
