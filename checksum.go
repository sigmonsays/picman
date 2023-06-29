package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

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
