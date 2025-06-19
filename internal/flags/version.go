package flags

import "github.com/urfave/cli/v2"

const (
	VersionKey = "version"
	BuildKey   = "build"
	CommitKey  = "commit"
	InfoKey    = "info"
)

var VersionFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    VersionKey,
		Usage:   "print a version tag or a short commit hash",
		Aliases: []string{"v"},
	},
	&cli.BoolFlag{
		Name:    BuildKey,
		Usage:   "prints short commit hash",
		Aliases: []string{CommitKey},
	},
	&cli.BoolFlag{
		Name:  InfoKey,
		Usage: "prints app info",
	},
}
