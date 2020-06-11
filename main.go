package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	cmd "github.com/shric/trpc/cmd/trpc"
)

type options struct {
	Common  cmd.CommonOptions `group:"global options"`
	Add     cmd.AddOptions    `command:"add" alias:"a" description:"Add torrents"`
	Errors  cmd.ErrorsOptions `command:"errors" alias:"e" description:"Show torrent error strings"`
	List    cmd.ListOptions   `command:"list" alias:"l" description:"List torrents"`
	Move    cmd.MoveOptions   `command:"move" alias:"mv" description:"Move torrent to another location"`
	Rm      cmd.RmOptions     `command:"rm" alias:"r" description:"Remove torrents"`
	Start   cmd.StartOptions  `command:"start" description:"Start torrents"`
	Stop    cmd.StopOptions   `command:"stop" description:"Start torrents"`
	Verify  cmd.VerifyOptions `command:"verify" alias:"hash" description:"Verify torrents (hash check)"`
	Version struct{}          `command:"version" description:"Print version"`
}

func main() {
	var args = new(options)

	p := flags.NewParser(args, flags.Default)
	remaining, err := p.Parse()

	if err != nil {
		os.Exit(1)
	}

	client := Connect()

	map[string]*cmd.Command{
		"verify":  cmd.NewCommand(cmd.Verify, args.Verify, remaining, args.Common, client),
		"rm":      cmd.NewCommand(cmd.Rm, args.Rm, remaining, args.Common, client),
		"start":   cmd.NewCommand(cmd.Start, args.Start, remaining, args.Common, client),
		"stop":    cmd.NewCommand(cmd.Stop, args.Stop, remaining, args.Common, client),
		"add":     cmd.NewCommand(cmd.Add, args.Add, remaining, args.Common, client),
		"list":    cmd.NewCommand(cmd.List, args.List, remaining, args.Common, client),
		"version": cmd.NewCommand(cmd.Version, args.Version, remaining, args.Common, client),
		"move":    cmd.NewCommand(cmd.Move, args.Move, remaining, args.Common, client),
	}[p.Active.Name].Run()
}
