package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

type Autosort struct {
	app *App
}

func (me *Autosort) Flags() []cli.Flag {
	incomingDir := "/data/Pictures-Android/AndroidDCIM/Camera"
	destDir := "/data/Pictures-incoming"
	ret := []cli.Flag{
		&cli.StringFlag{
			Name:    "src-dir",
			Usage:   "source directory",
			Aliases: []string{"s"},
			Value:   incomingDir,
		},
		&cli.StringFlag{
			Name:    "dest-dir",
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
			log.Warnf("ProcessFile: %s", path, err)
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

	// generate sha256
	sha256 := ""
	sha, err := Sha256File(fullpath)
	if err != nil {
		return err
	}
	sha256 = sha
	log.Tracef("Sha256 %s: %s", fullpath, sha256)

	return nil
}

var SupportedImageExt = map[string]bool{
	".hcif": true,
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

func Sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	cs := sha256.New()

	written, err := io.Copy(cs, f)
	if err != nil {
		return "", err
	}
	sha := cs.Sum(nil)
	shaStr := hex.EncodeToString(sha)
	log.Tracef("checksum %s: %s %d bytes", path, shaStr, written)

	return shaStr, nil
}
