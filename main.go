package main

import (
	"os"

	"github.com/sigmonsays/picman/autosort"
	"github.com/sigmonsays/picman/cleanup"
	"github.com/sigmonsays/picman/core"
	"github.com/sigmonsays/picman/list"
	"github.com/urfave/cli/v2"
)

func main() {

	appCtx := &core.App{}
	autosortCmd := &autosort.Autosort{appCtx}
	cleanupCmd := &cleanup.Cleanup{appCtx}
	listCmd := &list.List{appCtx}

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
			Action: autosortCmd.Action,
			Flags:  autosortCmd.Flags(),
		},
		{
			Name:   "cleanup",
			Usage:  "cleanup state files after import",
			Action: cleanupCmd.Action,
			Flags:  cleanupCmd.Flags(),
		},
		{
			Name:   "list",
			Usage:  "list all state files",
			Action: listCmd.Action,
			Flags:  listCmd.Flags(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Warnf("Picman: %s", err)
	}

}
