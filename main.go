package main

import (
	"github.com/helix-collective/hxgh/internal/gh"
	"github.com/helix-collective/hxgh/internal/types"

	"github.com/jpillora/opts"
)

var (
	version string = "dev"
	date    string = "na"
	commit  string = "na"
)

func main() {
	rootCmd := types.Root{}
	events := gh.NewEventsCmd(&rootCmd)
	e2csv := gh.NewCsvCmd(&rootCmd)
	opts.New(&rootCmd).
		Name("hxgh").
		EmbedGlobalFlagSet().
		Complete().
		Version(version).
		AddCommand(opts.New(events).Name("events").
			AddCommand(opts.New(e2csv).Name("csv"))).
		Parse().
		RunFatal()
}
