package cmd

import (
	"fmt"
	"math"
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/trpc/internal/util"

	"github.com/shric/trpc/internal/filter"
)

type setOptions struct {
	torrentOptions
	filter.Options `group:"filters"`
	ForceAll       bool  `long:"force-all" description:"Really limit all torrents"`
	Session        bool  `long:"session" short:"s" description:"Apply the limit to the session instead of torrent(s)"`
	DownLimit      int64 `long:"down" description:"Set download limit (0 for unlimited)" default:"9223372036854775807" default-mask:"-"`
	UpLimit        int64 `long:"up" description:"Set upload limit (0 for unlimited)" default:"9223372036854775807" default-mask:"-"`
}

// SessionLimit handles the limit command when --session is given.
func SessionLimit(c *Command) {
	opts, ok := c.Options.(setOptions)
	optionsCheck(ok)

	payload := &transmissionrpc.SessionArguments{}
	speedLimitDownEnabled := false
	speedLimitUpEnabled := false

	switch {
	case opts.DownLimit == math.MaxInt64:
	case opts.DownLimit > 0:
		speedLimitDownEnabled = true
		payload.SpeedLimitDown = &opts.DownLimit
		payload.SpeedLimitDownEnabled = &speedLimitDownEnabled

		c.statusf("Limited global download to %d KB/sec", opts.DownLimit)
	case opts.DownLimit <= 0:
		speedLimitDownEnabled = false
		payload.SpeedLimitDownEnabled = &speedLimitDownEnabled

		c.statusf("Removed global download limit")
	}

	switch {
	case opts.UpLimit == math.MaxInt64:
	case opts.UpLimit > 0:
		speedLimitUpEnabled = true
		payload.SpeedLimitUp = &opts.UpLimit
		payload.SpeedLimitUpEnabled = &speedLimitUpEnabled

		c.statusf("Limited global upload to %d KB/sec", opts.UpLimit)
	case opts.UpLimit <= 0:
		speedLimitUpEnabled = false
		payload.SpeedLimitUpEnabled = &speedLimitUpEnabled

		c.statusf("Removed global upload limit")
	}

	if !c.CommonOptions.DryRun {
		err := c.Client.SessionArgumentsSet(payload)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}

func setUploadLimit(payload *transmissionrpc.TorrentSetPayload, opts setOptions) string {
	uploadLimited := false
	switch {
	case opts.UpLimit == math.MaxInt64:
	case opts.UpLimit > 0:
		uploadLimited = true
		payload.UploadLimit = &opts.UpLimit
		payload.UploadLimited = &uploadLimited
		return fmt.Sprintf("Limiting upload to %d KB/sec for", opts.UpLimit)
	case opts.UpLimit <= 0:
		uploadLimited = false
		payload.UploadLimited = &uploadLimited
		return "Removed upload limit for"
	}
	return ""
}

func setDownloadLimit(payload *transmissionrpc.TorrentSetPayload, opts setOptions) string {
	downloadLimited := false

	switch {
	case opts.DownLimit == math.MaxInt64:
	case opts.DownLimit > 0:
		downloadLimited = true
		payload.DownloadLimit = &opts.DownLimit
		payload.DownloadLimited = &downloadLimited
		return fmt.Sprintf("Limited download to %d KB/sec for ", opts.DownLimit)
	case opts.DownLimit <= 0:
		downloadLimited = false
		payload.DownloadLimited = &downloadLimited
		return "Removed download limit for"
	}
	return ""
}

// TorrentLimit handles the limit command when --session isn't given.
func TorrentLimit(c *Command) {
	opts, ok := c.Options.(setOptions)
	optionsCheck(ok)
	util.ProcessTorrents(c.Client, opts.Options, opts.Positional.Torrents, commonArgs[:], func(torrent *transmissionrpc.Torrent) {
		IDs := make([]int64, 1)
		IDs[0] = *torrent.ID

		payload := &transmissionrpc.TorrentSetPayload{IDs: IDs}
		message := setDownloadLimit(payload, opts)
		c.status(message, torrent)

		message = setUploadLimit(payload, opts)
		c.status(message, torrent)

		if !c.CommonOptions.DryRun {
			err := c.Client.TorrentSet(payload)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
	}, nil, false)
}

// Set implements the limit command.
func Set(c *Command) {
	opts, ok := c.Options.(setOptions)
	optionsCheck(ok)

	if opts.UpLimit == math.MaxInt64 && opts.DownLimit == math.MaxInt64 {
		fmt.Fprint(os.Stderr, "Must specify either --down or --up\n")
		return
	}

	if len(opts.Positional.Torrents) == 0 && !opts.ForceAll && !opts.Session {
		fmt.Fprintln(os.Stderr,
			"Use --force-all if you really want to limit all torrents, use --session if you want to apply a session limit")
		return
	}

	if opts.Session && len(opts.Positional.Torrents) != 0 {
		fmt.Fprintln(os.Stderr,
			"Do not specify any torrents if using --session")
		return
	}

	if opts.Session {
		SessionLimit(c)
	} else {
		TorrentLimit(c)
	}
}
