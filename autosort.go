package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type Autosort struct {
	app *App
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

type AutoSortConfig struct {
	SourceDir string `yaml:"source_dir"`
	DestDir   string `yaml:"dest_dir"`
}

var conf = `
directories:
  - source_dir: /data/Pictures-Android/AndroidDCIM/Camera
    dest_dir: /data/Pictures-incoming/AndroidDCIM
`

func (me *Autosort) Action(c *cli.Context) error {
	incomingDir := c.String("src-dir")
	destDir := c.String("dest-dir")

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

func (me *Autosort) ProcessFile(root, fullpath string, info fs.FileInfo, dstdir string) error {
	if ShouldProcess(fullpath, info) == false {
		log.Tracef("ProcessFile %s: skipped", fullpath)
		return nil
	}

	relpath, err := filepath.Rel(root, fullpath)
	if err != nil {
		return err
	}

	srcSize := info.Size()

	// determine destination path name
	date, err := GetDate(fullpath)
	if err != nil {
		return err
	}

	year := date.Format("2006")
	month := date.Format("01")
	basename := filepath.Base(fullpath)
	destpath := filepath.Join(dstdir, year, month, basename)

	dstExists := false
	dstSize := int64(0)
	dst, err := os.Stat(destpath)
	if err == nil {
		dstExists = true
		dstSize = dst.Size()
	} else {
		dstExists = false
	}
	sizeMismatch := srcSize != dstSize

	if !dstExists {
		log.Tracef(" destination exists: %v", dstExists)
	}

	log.Tracef("src:%s size:%d mismatch:%v destpath:%s", relpath, srcSize, sizeMismatch, destpath)

	// do we need to convert it?
	// heic ->  jpg

	// generate sha256

	// sha256 := ""
	// sha, err := Sha256File(fullpath)
	// if err != nil {
	// 	return err
	// }
	// sha256 = sha
	// log.Tracef("Sha256 %s: %s", fullpath, sha256)

	return nil
}

var SupportedImageExt = map[string]bool{
	".heic": true,
	".jpeg": true,
	".jpg":  true,
	".png":  true,
	".mp4":  true,
}

func ShouldProcess(path string, info fs.FileInfo) bool {
	if info.Size() == 0 {
		log.Tracef("skip zero byte file %s", path)
		return false
	}

	base := filepath.Base(path)
	base = strings.ToLower(base)
	ext := filepath.Ext(base)

	supported, found := SupportedImageExt[ext]
	if found == false {
		log.Tracef("ext %s is not supported", ext)
		return false
	}
	if supported == false {
		return false
	}
	return true
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
