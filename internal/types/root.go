package types

// Root struct. Base cmd, with useful fields for subcommands
type Root struct {
	GhToken string `opts:"env=GITHUB_TOKEN"`
}
