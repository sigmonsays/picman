package core

import "fmt"

var SupportedImageExt = map[string]bool{
	".heic": true,
	".jpeg": true,
	".jpg":  true,
	".png":  true,
	".mov":  true,
	".mp4":  true,
}

func IsFileExtSupported(ext string) error {
	supported, found := SupportedImageExt[ext]
	if found == false {
		return fmt.Errorf("ext %s is not supported", ext)
	}
	if supported == false {
		return fmt.Errorf("Not supported")
	}
	return nil
}
