package core

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("core", func(newlog gologging.Logger) { log = newlog })
}
