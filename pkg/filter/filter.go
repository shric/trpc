package filter

import (
	"github.com/hekmon/transmissionrpc"
)

// Options declares all the command line arguments for filtering torrents
type Options struct {
	Incomplete bool `short:"i" long:"incomplete" description:"only incomplete torrents"`
}

type filterFunc struct {
	predicate func(*transmissionrpc.Torrent, string) bool
	set       interface{}
	args      []string
}

// Instance is used to hold all data required for a filter.
type Instance struct {
	opts        Options
	filterFuncs []filterFunc
	Args        []string
}

// New returns a new filter based on the options passed.
func New(opts Options) *Instance {
	filter := Instance{
		opts: opts,
		filterFuncs: []filterFunc{
			{
				predicate: func(t *transmissionrpc.Torrent, v string) bool {
					return *t.LeftUntilDone > 0
				},
				set:  &opts.Incomplete,
				args: []string{"leftUntilDone"},
			},
		},
		Args: make([]string, 0, 0),
	}
	argsSet := make(map[string]struct{})
	for _, i := range filter.filterFuncs {
		if set(i.set) {
			for _, arg := range i.args {
				argsSet[arg] = struct{}{}
			}
		}
	}
	for v := range argsSet {
		filter.Args = append(filter.Args, v)
	}
	return &filter
}

func set(set interface{}) bool {
	switch v := set.(type) {
	case *bool:
		return *v
	}
	return false
}

// CheckFilter checks if the supplied torrent matches after filters.
func (f *Instance) CheckFilter(torrent *transmissionrpc.Torrent) bool {
	match := true
	for _, fi := range f.filterFuncs {
		switch v := fi.set.(type) {
		case *bool:
			if *v && !fi.predicate(torrent, "") {
				match = false
			}
		}
	}
	return match
}
