package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shric/trpc/internal/fileutils"

	"github.com/hekmon/transmissionrpc"
)

const (
	fileIDunknown   = -2
	fileIDdirectory = -1
)

// Finder keeps all the state of a Finder instance returned by NewFinder.
type Finder struct {
	client   *transmissionrpc.Client
	torrents map[int64]*transmissionrpc.Torrent
	// First int64 is torrent ID, second int64 is file ID (-1 if a directory)
	cache           map[string][]int64
	incompleteDir   *string
	HasDownloadDirs bool
}

// NewFinder returns an instance of Finder.
func NewFinder(client *transmissionrpc.Client) *Finder {
	return &Finder{
		client:          client,
		torrents:        make(map[int64]*transmissionrpc.Torrent),
		cache:           make(map[string][]int64),
		incompleteDir:   getIncompleteDir(client),
		HasDownloadDirs: false,
	}
}

func (t *Finder) insertCache(path string, torrentID int64, fileID int64) {
	t.cache[path] = make([]int64, 2)
	t.cache[path][0] = torrentID
	t.cache[path][1] = fileID
}

func (t *Finder) getDownloadDirs() {
	// We only need to run this once.
	if t.HasDownloadDirs {
		return
	}

	torrents, err := t.client.TorrentGet([]string{"id", "downloadDir", "name"}, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var realDownloadDir string

	for _, torrent := range torrents {
		realDownloadDir = fileutils.RealPath(*torrent.DownloadDir)

		paths := []string{realDownloadDir}

		if t.incompleteDir != nil {
			paths = append(paths, *t.incompleteDir)
		}

		for _, path := range paths {
			fullPath := filepath.Join(path, *torrent.Name)
			eval, err := filepath.EvalSymlinks(fullPath)

			if err == nil {
				fullPath = eval
			}

			t.insertCache(fullPath, *torrent.ID, fileIDunknown)
			t.torrents[*torrent.ID] = torrent
		}
	}

	t.HasDownloadDirs = true
}

// Find returns the torrent and file ID of a given file.
func (t *Finder) Find(filename string) (*transmissionrpc.Torrent, int64) {
	absFilename := fileutils.RealPath(filename)

	if val, ok := t.cache[absFilename]; ok {
		if val[1] != fileIDunknown {
			return t.torrents[val[0]], val[1]
		}
	}

	t.getDownloadDirs()

	for fullPath, pair := range t.cache {
		if fullPath == absFilename && pair[1] != fileIDunknown {
			return t.torrents[pair[0]], pair[1]
		}

		if strings.HasPrefix(absFilename, fullPath) {
			if fileutils.IsDirectory(absFilename) {
				return t.torrents[t.cache[fullPath][0]], fileIDdirectory
			}

			torrents, err := t.client.TorrentGet([]string{"id", "downloadDir", "name", "files"}, []int64{pair[0]})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil, 0
			}

			torrent := torrents[0]
			realDownloadDir := fileutils.RealPath(*torrent.DownloadDir)

			for i, file := range torrent.Files {
				fullPath := filepath.Join(realDownloadDir, file.Name)
				t.insertCache(fullPath, *torrent.ID, int64(i))
			}

			return t.torrents[t.cache[absFilename][0]], t.cache[absFilename][1]
		}
	}

	return nil, 0
}
