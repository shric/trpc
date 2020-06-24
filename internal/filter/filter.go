package filter

import (
	"fmt"
	"net/url"
	"os"

	"github.com/shric/monkey/evaluator"
	"github.com/shric/monkey/lexer"
	"github.com/shric/monkey/object"
	"github.com/shric/monkey/parser"

	"github.com/hekmon/transmissionrpc"
)

// Options declares all the command line arguments for filtering torrents.
type Options struct {
	Filter     string `short:"f" long:"filter" description:"apply filter expression" unquote:"false"`
	Incomplete bool   `short:"i" long:"incomplete" description:"only incomplete torrents"`
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
		Args: make([]string, 0),
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
	default:
		fmt.Fprintln(os.Stderr, "Fatal internal error: unknown filter type")
		os.Exit(1)
	}

	return false
}

func envForTorrent(torrent *transmissionrpc.Torrent) *object.Environment {
	env := object.NewEnvironment()
	if *torrent.LeftUntilDone == 0 {
		env.Set("complete", evaluator.TRUE)
		env.Set("incomplete", evaluator.FALSE)
	} else {
		env.Set("complete", evaluator.FALSE)
		env.Set("incomplete", evaluator.TRUE)
	}
	env.Set("size", &object.Integer{Value: int64(torrent.SizeWhenDone.Byte())})
	trackers := make([]object.Object, len(torrent.Trackers))
	trackerStrings := make([]object.String, len(torrent.Trackers))
	for i, tracker := range torrent.Trackers {
		url, err := url.Parse(tracker.Announce)
		if err != nil {
			continue
		}
		trackerStrings[i] = object.String{Value: url.Hostname()}
		trackers[i] = &trackerStrings[i]
	}
	env.Set("trackers", &object.Array{Elements: trackers})
	return env
}

func checkFilterExpression(expr string, torrent *transmissionrpc.Torrent) bool {
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Fprintln(os.Stderr, "Filter parser error(s):")

		for _, msg := range p.Errors() {
			fmt.Fprintln(os.Stderr, "\t", msg)
		}

		fmt.Println(program.String())
		os.Exit(1)
	}

	result := evaluator.Eval(program, envForTorrent(torrent))
	switch v := result.(type) {
	case *object.Boolean:
		return v.Value
	case *object.Error:
		fmt.Fprintf(os.Stderr, "Invalid filter expression: %s\n", v.Message)
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "Invalid filter expression: doesn't evaluate to boolean: %q\n", v)
		os.Exit(1)
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

	if f.opts.Filter != "" {
		match = checkFilterExpression(f.opts.Filter, torrent)
	}

	return match
}
