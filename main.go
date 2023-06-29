package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	appCtx := &App{}
	autosort := &Autosort{appCtx}

	app := &cli.App{}
	app.Commands = []*cli.Command{
		{
			Name:   "autosort",
			Usage:  "autosort pictures",
			Action: autosort.Action,
			Flags:  autosort.Flags(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Warnf("Picman: %s", err)
	}

}
