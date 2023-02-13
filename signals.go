package byor

import (
	"os"
	"os/signal"
	"syscall"
)

func catchSignals() chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return sigs
}
