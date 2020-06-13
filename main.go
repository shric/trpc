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

	commandInstances := map[string]cmd.CommandInstance{
		"verify":  {Runner: cmd.Verify, Options: args.Verify},
		"errors":  {Runner: cmd.Errors, Options: args.Errors},
		"rm":      {Runner: cmd.Rm, Options: args.Rm},
		"start":   {Runner: cmd.Start, Options: args.Start},
		"stop":    {Runner: cmd.Stop, Options: args.Stop},
		"add":     {Runner: cmd.Add, Options: args.Add},
		"list":    {Runner: cmd.List, Options: args.List},
		"version": {Runner: cmd.Version, Options: args.Version},
		"move":    {Runner: cmd.Move, Options: args.Move},
	}

	command := &cmd.Command{
		PositionalArgs: remaining,
		CommonOptions:  args.Common,
		Client:         client,
	}

	command.CommandInstance = commandInstances[p.Active.Name]
	command.Run()
}
