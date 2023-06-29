package task

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("task", func(newlog gologging.Logger) { log = newlog })
}
