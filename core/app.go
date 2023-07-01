package core

import (
	"os"

	gologging "github.com/sigmonsays/go-logging"
)

type App struct {
}

func (me *App) Init() {
	gologging.SetLogOutput(os.Stdout)
}

func (me *App) SetLogLevel(lvl string) {
	gologging.SetLogLevel(lvl)
}
