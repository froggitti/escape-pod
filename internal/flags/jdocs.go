package flags

import "github.com/urfave/cli/v2"

var (
	JdocsDBName = "jdocs-db-name"
)

var JdocsFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    JdocsDBName,
		EnvVars: flagNamesToEnv(JdocsDBName),
	},
}
