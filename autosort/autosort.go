package autosort

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
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
		&cli.StringFlag{
			Name:    "output_file",
			Usage:   "write autosort results to json file",
			Aliases: []string{"o"},
		},
	}
	return ret
}

type Options struct {
	OneFile    string
	Force      bool
	NoCopy     bool
	OutputFile string
}

type AutosortResults struct {
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Stats       *Stats `json:"stats"`
	Rate        int    `json:"rate"`
	DurationSec int    `json:"duration_sec"`
}

func (me *Autosort) Action(c *cli.Context) error {
	sourceDir := c.String("source-dir")
	destDir := c.String("destination-dir")
	source := c.String("source")
	onefile := c.String("onefile")
	force := c.Bool("force")
	nocopy := c.Bool("no-copy")
	outputfile := c.String("output_file")

	opts := &Options{}
	opts.OneFile = onefile
	opts.Force = force
	opts.NoCopy = nocopy
	opts.OutputFile = outputfile

	err := me.PrepareSourceDir(sourceDir)
	if err != nil {
		return err
	}

	startTs := time.Now()
	stats := &Stats{}
	results := &AutosortResults{}
	results.Stats = stats
	results.StartTime = startTs.Format(time.RFC3339)

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
	log.Infof("source directory %s: processed %d files (%d copied) in %d ms %s",
		sourceDir, stats.Processed, stats.Copied, durMs, ratestr)

	results.EndTime = stopTs.Format(time.RFC3339)
	results.Rate = rate
	results.DurationSec = durSec

	resultsBuf, _ := json.MarshalIndent(results, "", " ")
	log.Tracef("results buffer:\n%s", resultsBuf)

	if opts.OutputFile != "" {
		err = ioutil.WriteFile(opts.OutputFile, resultsBuf, 0644)
		if err != nil {
			log.Warnf("WriteFile %s: %s", opts.OutputFile, err)
			return err
		}
	}

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

	stats.DirsProcessed++

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
