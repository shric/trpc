package cmd

import (
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/jessevdk/go-flags"
	"github.com/shric/trpc/internal/client"
)

// CommonOptions declares command line arguments that apply to all or most
// subcommands. It still needs to be explicitly included.
type commonOptions struct {
	DryRun bool `short:"n" long:"dry-run" description:"Dry run -- don't talk to the client, just print what would happen"`
	Debug  bool `short:"D" long:"debug" description:"Debug -- output the reply from server to stderr"`
}

// torrentOptions declares the positional command line argument for specifying 0 or more torrents.
type torrentOptions struct {
	Pos struct {
		Torrents []string `positional-arg-name:"torrent" description:"torrent ID or filename within torrent"`
	} `positional-args:"true"`
}

// fileOptions declares the positional command line argument for specifying 0 or more filenames.
type fileOptions struct {
	Pos struct {
		Files []string `positional-arg-name:"file" description:"filename"`
	} `positional-args:"true"`
}

type options struct {
	Common  commonOptions `group:"global options"`
	Add     addOptions    `command:"add" alias:"a" description:"Add torrents"`
	Errors  errorsOptions `command:"errors" alias:"e" description:"Show torrent error strings"`
	List    listOptions   `command:"list" alias:"l" description:"List torrents"`
	Move    moveOptions   `command:"move" alias:"mv" description:"Move torrent to another location"`
	Rename  renameOptions `command:"rename" description:"Rename torrent file"`
	Rm      rmOptions     `command:"rm" alias:"r" description:"Remove torrents"`
	Set     setOptions    `command:"set" description:"Set torrent priorities/speeds or session speeds"`
	Start   startOptions  `command:"start" description:"Start torrents"`
	Stop    stopOptions   `command:"stop" description:"Start torrents"`
	Verify  verifyOptions `command:"verify" alias:"hash" description:"Verify torrents (hash check)"`
	Watch   watchOptions  `command:"watch" description:"Watch progress for torrents"`
	Which   whichOptions  `command:"which" description:"Identify which file/path a torrent belongs to"`
	Version struct{}      `command:"version" description:"Print version"`
}

// CommandInstance is the data specific to one command.
type CommandInstance struct {
	Options interface{}
	Runner  func(c *Command)
}

// Command holds everything needed to run a command.
type Command struct {
	PositionalArgs []string
	CommonOptions  commonOptions
	Client         *transmissionrpc.Client
	CommandInstance
}

var commonArgs = [...]string{
	"id", "name", "status", "name", "recheckProgress", "sizeWhenDone",
	"rateUpload", "eta", "id", "leftUntilDone", "recheckProgress", "error",
	"rateDownload", "status", "trackers", "bandwidthPriority", "uploadedEver",
	"downloadDir", "addedDate", "doneDate", "startDate", "isFinished",
	"errorString",
}

// Run parses flags.
func Run() {
	var args = new(options)

	p := flags.NewParser(args, flags.Default)
	_, err := p.Parse()

	if err != nil {
		os.Exit(1)
	}

	commandInstances := map[string]CommandInstance{
		"add":     {Runner: Add, Options: args.Add},
		"errors":  {Runner: Errors, Options: args.Errors},
		"list":    {Runner: List, Options: args.List},
		"move":    {Runner: Move, Options: args.Move},
		"rename":  {Runner: Rename, Options: args.Rename},
		"rm":      {Runner: Rm, Options: args.Rm},
		"set":     {Runner: Set, Options: args.Set},
		"start":   {Runner: Start, Options: args.Start},
		"stop":    {Runner: Stop, Options: args.Stop},
		"verify":  {Runner: Verify, Options: args.Verify},
		"version": {Runner: Version, Options: args.Version},
		"watch":   {Runner: Watch, Options: args.Watch},
		"which":   {Runner: Which, Options: args.Which},
	}

	c := client.Connect(args.Common.Debug)

	command := &Command{
		CommonOptions: args.Common,
		Client:        c,
	}

	command.CommandInstance = commandInstances[p.Active.Name]
	command.Run()
}

// Run is a simple wrapper to call the runner function of a command.
func (c *Command) Run() {
	if c.Runner != nil {
		c.Runner(c)
	} else {
		fmt.Fprintln(os.Stderr, "Fatal internal error: command not implemented")
	}
}

func (c *Command) statusf(format string, a ...interface{}) {
	var dryRun string
	if c.CommonOptions.DryRun {
		dryRun = "[dry run] "
	}

	fmt.Printf(dryRun+format+"\n", a...)
}

func (c *Command) status(msg string, torrent *transmissionrpc.Torrent) {
	c.statusf("%s %d: %s", msg, *torrent.ID, *torrent.Name)
}

func optionsCheck(ok bool) {
	if !ok {
		fmt.Fprintln(os.Stderr, "Fatal internal error: bad options passed.")
		os.Exit(1)
	}
}

// Constants for binary units.
const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB
)
