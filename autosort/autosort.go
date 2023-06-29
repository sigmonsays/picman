package autosort

import (
	"bytes"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
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
	}
	return ret
}

func (me *Autosort) Action(c *cli.Context) error {
	incomingDir := c.String("source-dir")
	destDir := c.String("destination-dir")

	// copy files from incoming directories (and keep track of them)
	// once copied to an incoming directory we mark them in the DB

	err := me.ProcessDir(incomingDir, destDir)
	if err != nil {
		return err
	}

	return nil
}

func (me *Autosort) ProcessDir(srcdir, dstdir string) error {
	log.Tracef("ProcessDir %s", srcdir)

	walkfn := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Warnf("walk: %s: %s", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		err = me.ProcessFile(srcdir, path, info, dstdir)
		if err != nil {
			log.Warnf("ProcessFile %s: %s", path, err)
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
