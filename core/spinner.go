package core

import (
	"time"

	"github.com/briandowns/spinner"
)

var spin = spinner.New(spinner.CharSets[39], 100*time.Millisecond)

func StartSpinner(message string) {
	if spin == nil || spin.Active() {
		return
	}

	spin.Suffix = " " + message
	spin.Start()
}

func StopSpinner(message string) {
	if spin == nil || !spin.Active() {
		return
	}

	spin.FinalMSG = "✅ " + message + "\n"
	spin.Stop()
}
