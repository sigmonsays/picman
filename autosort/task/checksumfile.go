package task

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/sigmonsays/picman/core"
)

// // determine destination path name
// date, err := GetDate(fullpath)
// if err != nil {
// 	return err
// }

func NewChecksumFile(w *core.Workflow) *ChecksumFile {
	ret := &ChecksumFile{}
	ret.Workflow = w
	return ret
}

type ChecksumFile struct {
	Workflow *core.Workflow
}

func (me *ChecksumFile) Run(state *core.State) error {
	log.Tracef("start %s", state.OriginalFilename)
	if state.Checksum.Sha256 != "" {
		log.Tracef("Checksum already generated, skipping")
		return nil
	}

	cs, err := Sha256File(state.OriginalFilename)
	if err != nil {
		return state.StopProcessing("Sha256 %s", err)
	}

	state.Checksum.Sha256 = cs

	return nil
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
