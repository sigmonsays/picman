package main

import (
	"os"

	"github.com/sigmonsays/picman/autosort"
	"github.com/sigmonsays/picman/core"
	"github.com/urfave/cli/v2"
)

func main() {

	appCtx := &core.App{}
	autosort := &autosort.Autosort{appCtx}

	app := &cli.App{}
	appCtx.Init()
	app.Before = func(c *cli.Context) error {
		loglevel := c.String("loglevel")
		if loglevel != "" {
			appCtx.SetLogLevel(loglevel)
		}
		return nil
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "loglevel",
			Usage:   "set log level",
			Aliases: []string{"l"},
			Value:   "INFO",
		},
	}
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
