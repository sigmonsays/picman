package list

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sigmonsays/picman/core"
	"github.com/urfave/cli/v2"
)

type List struct {
	App *core.App
}

func (me *List) Flags() []cli.Flag {
	// incomingDir := "/data/Pictures-Android/AndroidDCIM/Camera"
	incomingDir, err := os.Getwd()
	if err != nil {
		log.Warnf("Getwd %s", err)
	}

	ret := []cli.Flag{
		&cli.StringFlag{
			Name:    "source-dir",
			Usage:   "source directory",
			Aliases: []string{"s"},
			Value:   incomingDir,
		},
		&cli.StringFlag{
			Name:    "onefile",
			Usage:   "process just one file",
			Aliases: []string{""},
		},
	}
	return ret
}

type Options struct {
	OneFile string
}

func (me *List) Action(c *cli.Context) error {
	sourceDir := c.String("source-dir")
	onefile := c.String("onefile")

	opts := &Options{}
	opts.OneFile = onefile

	var err error

	err = me.ProcessDir(sourceDir, opts)
	if err != nil {
		return err
	}
	return nil
}

func (me *List) ProcessDir(srcdir string, opts *Options) error {
	log.Tracef("ProcessDir %s", srcdir)

	walkfn := func(statefile string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Warnf("walk: %s: %s", statefile, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		err = me.ProcessFile(srcdir, statefile, opts)
		if err != nil {
			log.Warnf("ProcessFile %s: %s", statefile, err)
		}
		return nil
	}
	statedir := filepath.Join(srcdir, core.StateSubDir)
	err := filepath.Walk(statedir, walkfn)
	if err != nil {
		return err
	}
	return nil
}
