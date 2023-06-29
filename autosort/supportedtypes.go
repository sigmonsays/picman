package autosort

import (
	"io/fs"
	"path/filepath"
	"strings"
)

var SupportedImageExt = map[string]bool{
	".heic": true,
	".jpeg": true,
	".jpg":  true,
	".png":  true,
	".mov":  true,
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
