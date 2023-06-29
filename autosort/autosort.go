package autosort

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sigmonsays/picman/core"
	"github.com/urfave/cli/v2"
)

var DirMask = os.FileMode(0755)

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
			Value:   destDir,
		},
	}
	return ret
}

type Options struct {
	OneFile string
	Force   bool
}

func (me *Autosort) Action(c *cli.Context) error {
	sourceDir := c.String("source-dir")
	destDir := c.String("destination-dir")
	source := c.String("source")
	onefile := c.String("onefile")
	force := c.Bool("force")

	opts := &Options{}
	opts.OneFile = onefile
	opts.Force = force

	err := me.PrepareSourceDir(sourceDir)
	if err != nil {
		return err
	}

	if onefile != "" {

		fullpath := onefile
		if x, err := filepath.Abs(onefile); err == nil {
			fullpath = x
		}
		info, err := os.Stat(fullpath)
		if err != nil {
			return err
		}

		err = me.ProcessFile(sourceDir, fullpath, info, destDir, source, opts)
		if err != nil {
			return err
		}
		return nil
	}

	err = me.ProcessDir(sourceDir, destDir, source, opts)
	if err != nil {
		return err
	}

	return nil
}

func (me *Autosort) PrepareSourceDir(srcdir string) error {
	// begin procesing
	statedir := filepath.Join(srcdir, StateSubDir)
	os.MkdirAll(statedir, DirMask)
	errordir := filepath.Join(srcdir, ErrorSubDir)
	os.MkdirAll(errordir, DirMask)

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

func (me *Autosort) ProcessDir(srcdir, dstdir string, source string, opts *Options) error {
	log.Tracef("ProcessDir %s", srcdir)

	walkfn := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Warnf("walk: %s: %s", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		err = me.ProcessFile(srcdir, path, info, dstdir, source, opts)
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

func GetDate(path string) (time.Time, error) {
	args := []string{
		"exiftool",
		"-S",
		"-DateTimeOriginal",
		path,
	}
	none := time.Time{}
	buf := bytes.NewBuffer(nil)
	c := exec.Command(args[0], args[1:]...)
	c.Stdout = buf
	err := c.Run()
	if err != nil {
		return none, err
	}
	tmp := strings.Trim(buf.String(), "\n")
	vals := strings.SplitN(tmp, ":", 2)
	if len(vals) < 2 {
		return none, fmt.Errorf("short length")
	}
	t := strings.Trim(vals[1], "\n\t ")

	// 2023:03:02 18:13:31
	tm, err := time.Parse("2006:01:02 15:04:05", t)
	if len(vals) < 2 {
		return none, fmt.Errorf("parse time %s: %s", t, err)
	}

	return tm, nil
}
