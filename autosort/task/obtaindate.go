package task

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// // determine destination path name
// date, err := GetDate(fullpath)
// if err != nil {
// 	return err
// }

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
