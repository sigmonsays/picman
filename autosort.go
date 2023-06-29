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
	ret := []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "config",
			Usage:   "config file",
			Aliases: []string{"c"},
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

func (me *Autosort) Action(ctx *cli.Context) error {

	// copy files from incoming directories (and keep track of them)
	// once copied to an incoming directory we mark them in the DB
	incomingDir := "/data/Pictures-Android/AndroidDCIM/Camera"

	err = me.ProcessFiles(picDb, incomingDir, destDir)
	if err != nil {
		return err
	}

	return nil
}

func (me *Autosort) CopyIncomingFiles(picDb *db.Db, dir string) error {
	log.Tracef("CopyIncomingFiles %s", dir)

	walkfn := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Warnf("walk: %s: %s", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		err = me.ProcessFile(picDb, dir, path, info)
		if err != nil {
			log.Warnf("ProcessFile: %s", dir, err)
		}
		return nil
	}
	err := filepath.Walk(dir, walkfn)
	if err != nil {
		return err
	}
	return nil
}

func (me *Autosort) ProcessFile(picDb *db.Db, root, fullpath string, info fs.FileInfo) error {

	assetFound := false
	sha256 := ""
	row, err := picDb.GetAssetRow(fullpath)
	if err == nil {
		assetFound = true
	}
	if err != nil {
		log.Tracef("GetAssetRow %s: %s", fullpath, err)
	}

	log.Tracef("ProcessFile %s: found:%v", fullpath, assetFound)
	if row != nil {
		log.Tracef("sha256 %s", row.Sha256)
	}
	if ShouldProcess(fullpath, info) == false {
		log.Tracef("ProcessFile %s: found:%v skipped", fullpath, assetFound)
		return nil
	}

	if assetFound == false {
		sha, err := Sha256File(fullpath)
		if err != nil {
			return err
		}
		log.Tracef("Sha256 %s: %s", fullpath, sha)
		sha256 = sha
		row := &db.AssetRow{}
		row.FullPath = fullpath
		row.Sha256 = sha256
		err = picDb.InsertAssetRow(row)
		if err != nil {
			return err
		}
	}

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
