package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/filter"
)

// Command holds everything needed to run a command.
type Command struct {
	CommandOptions interface{}
	PositionalArgs []string
	CommonOptions  CommonOptions
	Client         *transmissionrpc.Client
	Runner         func(c *Command)
}

// Run is a simple wrapper to call the runner function of a command.
func (c *Command) Run() {
	if c != nil {
		c.Runner(c)
	}
}

// NewCommand returns a Command.
func NewCommand(runner func(c *Command), commandOptions interface{},
	positionalArgs []string, commonOptions CommonOptions,
	client *transmissionrpc.Client) (command *Command) {
	return &Command{
		CommandOptions: commandOptions,
		PositionalArgs: positionalArgs,
		CommonOptions:  commonOptions,
		Client:         client,
		Runner:         runner,
	}
}

// CommonOptions declares command line arguments that apply to all or most
// subcommands. It still needs to be explicitly included.
type CommonOptions struct {
	DryRun bool `short:"n" long:"dry-run" description:"Dry run -- don't talk to the client, just print what would happen"`
}

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

	enabled = true
	incompleteDir = *session.IncompleteDir

	return
}

// getids attempts to convert a list of torrent filenames to their corresponding ID
// numbers in transmission.
func getids(client *transmissionrpc.Client, fnames []string) []int64 {
	// Let's do no work if given an empty list as this function is expensive
	if len(fnames) == 0 {
		return nil
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

// ProcessTorrents runs the supplied function over all torrents matching the args and filters.
func ProcessTorrents(client *transmissionrpc.Client, filterOptions filter.Options, args []string,
	fields []string, do func(torrent *transmissionrpc.Torrent),
) {
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

func (c *Command) status(msg string, torrent *transmissionrpc.Torrent) {
	var dryRun string
	if c.CommonOptions.DryRun {
		dryRun = "[dry run] "
	}

	fmt.Printf("%s%s %d: %s\n", dryRun, msg, *torrent.ID, *torrent.Name)
}

func optionsCheck(ok bool) {
	if !ok {
		fmt.Fprintln(os.Stderr, "Fatal internal error: bad options passed.")
		os.Exit(1)
	}
}
