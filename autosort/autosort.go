package autosort

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/sigmonsays/picman/core"
	"github.com/urfave/cli/v2"
)

type Autosort struct {
	App *core.App
}

func (me *Autosort) Flags() []cli.Flag {
	incomingDir := "/data/Pictures-Android/AndroidDCIM/Camera"
	destDir := "/data/Pictures"
	ret := []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Usage:   "start without previous state",
			Aliases: []string{"f"},
		},
		&cli.BoolFlag{
			Name:  "no-copy",
			Usage: "do not copy to final destination",
		},
		&cli.StringFlag{
			Name:    "source",
			Usage:   "source",
			Aliases: []string{"S"},
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

func (me *Autosort) Action(c *cli.Context) error {
	sourceDir := c.String("source-dir")
	destDir := c.String("destination-dir")
	source := c.String("source")
	onefile := c.String("onefile")
	force := c.Bool("force")
	nocopy := c.Bool("no-copy")

	opts := &Options{}
	opts.OneFile = onefile
	opts.Force = force
	opts.NoCopy = nocopy

	err := me.PrepareSourceDir(sourceDir)
	if err != nil {
		return err
	}

	startTs := time.Now()
	stats := &Stats{}

	if onefile != "" {
		fullpath := onefile
		if x, err := filepath.Abs(onefile); err == nil {
			fullpath = x
		}
		log.Tracef("onefile test %s", fullpath)
		info, err := os.Stat(fullpath)
		if err != nil {
			return err
		}

		err = me.ProcessFile(sourceDir, fullpath, info, destDir, source, opts, stats)
		if err != nil {
			return err
		}
		return nil
	}

	err = me.ProcessDir(sourceDir, destDir, source, opts, stats)
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

func (me *Autosort) PrepareSourceDir(srcdir string) error {
	// begin procesing
	statedir := filepath.Join(srcdir, StateSubDir)
	errordir := filepath.Join(srcdir, ErrorSubDir)

	for i := 0; i <= 255; i++ {
		h := fmt.Sprintf("%02x", i)
		d := filepath.Join(statedir, h)
		os.MkdirAll(d, core.DirMask)
		d = filepath.Join(errordir, h)
		os.MkdirAll(d, core.DirMask)
	}

	// ensure statedir exists
	st, err := os.Stat(statedir)
	if err != nil {
		return fmt.Errorf("state dir %s does not exist: %s", statedir, err)
	}
	if st.IsDir() == false {
		return fmt.Errorf("state dir %s: is not a directory", statedir)
	}

	return nil
}

func (me *Autosort) ProcessDir(srcdir, dstdir string, source string, opts *Options, stats *Stats) error {
	log.Tracef("ProcessDir %s", srcdir)

	walkfn := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Warnf("walk: %s: %s", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		err = me.ProcessFile(srcdir, path, info, dstdir, source, opts, stats)
		if err != nil {
			if err == core.StopWorkflow {
				return fmt.Errorf("%s indicates stop workflow", path)
			} else if err == core.StopProcessing {
				// stop processing current file and advance to next
				log.Tracef("%s stop processing file", path)
				return nil
			}
			log.Warnf("ProcessFile %s: %s", path, err)
			return nil
		}
		return nil
	}
	err := filepath.Walk(srcdir, walkfn)
	if err != nil {
		return err
	}
	return nil
}
