package cleanup

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/sigmonsays/picman/core"
	"github.com/urfave/cli/v2"
)

type Cleanup struct {
	App *core.App
}

func (me *Cleanup) Flags() []cli.Flag {
	incomingDir := "/data/Pictures-Android/AndroidDCIM/Camera"
	destDir := "/data/Pictures"

	ret := []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Usage:   "start without previous state",
			Aliases: []string{"f"},
		},
		&cli.StringFlag{
			Name:    "source-dir",
			Usage:   "source directory",
			Aliases: []string{"s"},
			Value:   incomingDir,
		},
		&cli.StringFlag{
			Name:    "destination-dir",
			Usage:   "destination directory",
			Aliases: []string{"d"},
			Value:   destDir,
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
	Force   bool
	NoCopy  bool
}

func (me *Cleanup) Action(c *cli.Context) error {
	sourceDir := c.String("source-dir")
	onefile := c.String("onefile")
	force := c.Bool("force")

	opts := &Options{}
	opts.OneFile = onefile
	opts.Force = force

	startTs := time.Now()
	stats := &Stats{}

	var err error

	if onefile != "" {
		statefile := onefile
		if x, err := filepath.Abs(onefile); err == nil {
			statefile = x
		}
		log.Tracef("onefile test %s", statefile)

		err = me.ProcessFile(sourceDir, statefile, opts, stats)
		if err != nil {
			return err
		}
		return nil
	}

	err = me.ProcessDir(sourceDir, opts, stats)
	if err != nil {
		return err
	}
	stopTs := time.Now()
	dur := stopTs.Sub(startTs)
	durMs := int64(dur.Milliseconds())
	rate := 0
	durSec := int(dur.Seconds())
	ratestr := ""
	if durSec > 0 && stats.Processed > 0 {
		rate = stats.Processed / durSec
		ratestr = fmt.Sprintf("(%d files/sec)", rate)
	}
	log.Infof("processed %d files in %d ms %s", stats.Processed, durMs, ratestr)
	return nil
}

func (me *Cleanup) ProcessDir(srcdir string, opts *Options, stats *Stats) error {
	log.Tracef("ProcessDir %s", srcdir)

	walkfn := func(statefile string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Warnf("walk: %s: %s", statefile, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		err = me.ProcessFile(srcdir, statefile, opts, stats)
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
